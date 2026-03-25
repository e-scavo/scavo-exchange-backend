package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/e-scavo/scavo-exchange-backend/internal/thirdparty/sha3local"
)

var (
	ErrInvalidWalletSignature = errors.New("invalid wallet signature")
)

var secp256k1 elliptic.Curve = newSecp256k1Curve()

type secp256k1Curve struct {
	params *elliptic.CurveParams
}

func newSecp256k1Curve() elliptic.Curve {
	p := &elliptic.CurveParams{Name: "secp256k1"}
	p.P = mustHexBig("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
	p.N = mustHexBig("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")
	p.B = big.NewInt(7)
	p.Gx = mustHexBig("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798")
	p.Gy = mustHexBig("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8")
	p.BitSize = 256
	return &secp256k1Curve{params: p}
}

func (c *secp256k1Curve) Params() *elliptic.CurveParams { return c.params }

func (c *secp256k1Curve) IsOnCurve(x, y *big.Int) bool {
	if isInfinity(x, y) {
		return false
	}
	if x.Sign() < 0 || x.Cmp(c.params.P) >= 0 || y.Sign() < 0 || y.Cmp(c.params.P) >= 0 {
		return false
	}
	left := new(big.Int).Mul(y, y)
	left.Mod(left, c.params.P)
	right := new(big.Int).Exp(x, big.NewInt(3), c.params.P)
	right.Add(right, c.params.B)
	right.Mod(right, c.params.P)
	return left.Cmp(right) == 0
}

func (c *secp256k1Curve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	if isInfinity(x1, y1) {
		return clonePoint(x2, y2)
	}
	if isInfinity(x2, y2) {
		return clonePoint(x1, y1)
	}
	if !c.IsOnCurve(x1, y1) || !c.IsOnCurve(x2, y2) {
		panic("point not on secp256k1 curve")
	}

	p := c.params.P
	if x1.Cmp(x2) == 0 {
		ysum := new(big.Int).Add(y1, y2)
		ysum.Mod(ysum, p)
		if ysum.Sign() == 0 {
			return big.NewInt(0), big.NewInt(0)
		}
		return c.Double(x1, y1)
	}

	num := new(big.Int).Sub(y2, y1)
	num.Mod(num, p)
	den := new(big.Int).Sub(x2, x1)
	den.Mod(den, p)
	denInv := new(big.Int).ModInverse(den, p)
	if denInv == nil {
		return big.NewInt(0), big.NewInt(0)
	}
	lambda := new(big.Int).Mul(num, denInv)
	lambda.Mod(lambda, p)

	x3 := new(big.Int).Mul(lambda, lambda)
	x3.Sub(x3, x1)
	x3.Sub(x3, x2)
	x3.Mod(x3, p)
	if x3.Sign() < 0 {
		x3.Add(x3, p)
	}

	y3 := new(big.Int).Sub(x1, x3)
	y3.Mul(lambda, y3)
	y3.Sub(y3, y1)
	y3.Mod(y3, p)
	if y3.Sign() < 0 {
		y3.Add(y3, p)
	}

	return x3, y3
}

func (c *secp256k1Curve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	if isInfinity(x1, y1) {
		return big.NewInt(0), big.NewInt(0)
	}
	if !c.IsOnCurve(x1, y1) {
		panic("point not on secp256k1 curve")
	}
	if y1.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0)
	}

	p := c.params.P
	num := new(big.Int).Mul(x1, x1)
	num.Mul(num, big.NewInt(3))
	num.Mod(num, p)
	den := new(big.Int).Mul(y1, big.NewInt(2))
	den.Mod(den, p)
	denInv := new(big.Int).ModInverse(den, p)
	if denInv == nil {
		return big.NewInt(0), big.NewInt(0)
	}
	lambda := new(big.Int).Mul(num, denInv)
	lambda.Mod(lambda, p)

	x3 := new(big.Int).Mul(lambda, lambda)
	twoX1 := new(big.Int).Mul(x1, big.NewInt(2))
	x3.Sub(x3, twoX1)
	x3.Mod(x3, p)
	if x3.Sign() < 0 {
		x3.Add(x3, p)
	}

	y3 := new(big.Int).Sub(x1, x3)
	y3.Mul(lambda, y3)
	y3.Sub(y3, y1)
	y3.Mod(y3, p)
	if y3.Sign() < 0 {
		y3.Add(y3, p)
	}

	return x3, y3
}

func (c *secp256k1Curve) ScalarMult(Bx, By *big.Int, k []byte) (*big.Int, *big.Int) {
	if isInfinity(Bx, By) {
		return big.NewInt(0), big.NewInt(0)
	}
	if !c.IsOnCurve(Bx, By) {
		panic("point not on secp256k1 curve")
	}

	var x, y *big.Int
	for _, b := range k {
		for bit := 7; bit >= 0; bit-- {
			x, y = c.Double(x, y)
			if (b>>uint(bit))&1 == 1 {
				x, y = c.Add(x, y, Bx, By)
			}
		}
	}
	if isInfinity(x, y) {
		return big.NewInt(0), big.NewInt(0)
	}
	return x, y
}

func (c *secp256k1Curve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	return c.ScalarMult(c.params.Gx, c.params.Gy, k)
}

func isInfinity(x, y *big.Int) bool {
	return x == nil || y == nil || (x.Sign() == 0 && y.Sign() == 0)
}

func clonePoint(x, y *big.Int) (*big.Int, *big.Int) {
	if isInfinity(x, y) {
		return big.NewInt(0), big.NewInt(0)
	}
	return new(big.Int).Set(x), new(big.Int).Set(y)
}

func mustHexBig(v string) *big.Int {
	n := new(big.Int)
	if _, ok := n.SetString(v, 16); !ok {
		panic("invalid hex big int")
	}
	return n
}

func ethereumMessageHash(message string) []byte {
	payload := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	h := sha3local.NewLegacyKeccak256()
	_, _ = h.Write([]byte(payload))
	return h.Sum(nil)
}

func recoverWalletAddress(message, signature string) (string, error) {
	sig, err := decodeSignature(signature)
	if err != nil {
		return "", err
	}

	hash := ethereumMessageHash(message)
	x, y, err := recoverSecp256k1Pubkey(hash, sig)
	if err != nil {
		return "", err
	}

	return publicKeyToAddress(x, y), nil
}

func decodeSignature(signature string) ([]byte, error) {
	signature = strings.TrimSpace(signature)
	signature = strings.TrimPrefix(strings.ToLower(signature), "0x")
	if len(signature) != 130 {
		return nil, ErrInvalidWalletSignature
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return nil, ErrInvalidWalletSignature
	}
	if len(sig) != 65 {
		return nil, ErrInvalidWalletSignature
	}

	v := sig[64]
	switch {
	case v >= 27 && v <= 30:
		sig[64] = v - 27
	case v <= 3:
		// already normalized
	default:
		return nil, ErrInvalidWalletSignature
	}

	return sig, nil
}

func recoverSecp256k1Pubkey(hash []byte, sig []byte) (*big.Int, *big.Int, error) {
	if len(hash) == 0 || len(sig) != 65 {
		return nil, nil, ErrInvalidWalletSignature
	}

	curve := secp256k1
	params := curve.Params()
	n := params.N
	p := params.P

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := int(sig[64])

	if r.Sign() <= 0 || s.Sign() <= 0 || r.Cmp(n) >= 0 || s.Cmp(n) >= 0 {
		return nil, nil, ErrInvalidWalletSignature
	}
	if v < 0 || v > 3 {
		return nil, nil, ErrInvalidWalletSignature
	}

	j := v / 2
	x := new(big.Int).Mul(big.NewInt(int64(j)), n)
	x.Add(x, r)
	if x.Cmp(p) >= 0 {
		return nil, nil, ErrInvalidWalletSignature
	}

	y, err := secp256k1Y(x, v%2)
	if err != nil {
		return nil, nil, err
	}
	if !curve.IsOnCurve(x, y) {
		return nil, nil, ErrInvalidWalletSignature
	}

	rInv := new(big.Int).ModInverse(r, n)
	if rInv == nil {
		return nil, nil, ErrInvalidWalletSignature
	}

	e := new(big.Int).SetBytes(hash)
	e.Mod(e, n)
	eNeg := new(big.Int).Neg(e)
	eNeg.Mod(eNeg, n)

	sx, sy := curve.ScalarMult(x, y, s.Bytes())
	ex, ey := curve.ScalarBaseMult(eNeg.Bytes())
	qx, qy := curve.Add(sx, sy, ex, ey)
	if isInfinity(qx, qy) {
		return nil, nil, ErrInvalidWalletSignature
	}

	qx, qy = curve.ScalarMult(qx, qy, rInv.Bytes())
	if isInfinity(qx, qy) || !curve.IsOnCurve(qx, qy) {
		return nil, nil, ErrInvalidWalletSignature
	}

	pub := ecdsa.PublicKey{Curve: curve, X: qx, Y: qy}
	if !ecdsa.Verify(&pub, hash, r, s) {
		return nil, nil, ErrInvalidWalletSignature
	}

	return qx, qy, nil
}

func secp256k1Y(x *big.Int, parity int) (*big.Int, error) {
	p := secp256k1.Params().P
	y2 := new(big.Int).Exp(x, big.NewInt(3), p)
	y2.Add(y2, big.NewInt(7))
	y2.Mod(y2, p)

	exponent := new(big.Int).Add(p, big.NewInt(1))
	exponent.Rsh(exponent, 2)
	y := new(big.Int).Exp(y2, exponent, p)

	check := new(big.Int).Mul(y, y)
	check.Mod(check, p)
	if check.Cmp(y2) != 0 {
		return nil, ErrInvalidWalletSignature
	}

	if y.Bit(0) != uint(parity&1) {
		y.Sub(p, y)
		y.Mod(y, p)
	}

	return y, nil
}

func publicKeyToAddress(x, y *big.Int) string {
	if isInfinity(x, y) {
		return ""
	}

	buf := make([]byte, 64)
	x.FillBytes(buf[:32])
	y.FillBytes(buf[32:])

	h := sha3local.NewLegacyKeccak256()
	_, _ = h.Write(buf)
	sum := h.Sum(nil)
	return "0x" + hex.EncodeToString(sum[len(sum)-20:])
}

func walletUserID(address string) string {
	return "wallet:" + normalizeWalletAddress(address)
}

func normalizeWalletAddress(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}
