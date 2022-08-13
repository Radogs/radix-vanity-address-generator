package wallet

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip32"
)

func TestSeedFromMnemonic(t *testing.T) {
	testCases := []struct {
		name            string
		mnemonic        string
		expectedRootKey string
	}{
		{
			name:            "Should derive correct root key from mnemonic",
			mnemonic:        "special obey degree parrot lumber funny celery resist broken bundle robust popular",
			expectedRootKey: "xprv9s21ZrQH143K3pkKxfdeuMvyLKie9WgH7X7FdtiQEp1URb1qv4S71SyvxtE2NPqb9dsnk8UmhJ7mAexd9wo95cTCy4sgXGat7PfA87fUwQm",
		},

		{
			name:            "Should derive correct root key from different mnemonic",
			mnemonic:        "hospital claim choose expose artwork camera trumpet seminar junk ship client ocean",
			expectedRootKey: "xprv9s21ZrQH143K3fvGTe2ooHnPbaqC67dxfNK2KWYQTcBgAydTAo85tFyLwDCbgcdNBDE5DwCYovnaubVewoRii6ug5kUbqURb76xUNQ9Gbss",
		},
	}

	for _, testCase := range testCases {
		seed := seedFromMnemonic(testCase.mnemonic)
		key, err := bip32.NewMasterKey(seed)
		assert.Nil(t, err)
		assert.Equal(t, testCase.expectedRootKey, key.String(), testCase.name)
	}
}

func TestDerivationPathFromMnemonic(t *testing.T) {
	testCases := []struct {
		name               string
		mnemonic           string
		extendedPrivateKey string
		extendedPublicKey  string
	}{
		{
			name:               "Should derive correct root keys from mnemonic",
			mnemonic:           "special obey degree parrot lumber funny celery resist broken bundle robust popular",
			extendedPrivateKey: "xprvA1tW2qA4Mxf7ELMG6YtiSbiFK3mYhVbEnZMtfnUUXri7bxdnYpWUBWWGPcHTJe4YjCHHUNVS9WEwhrh2gVEVm1enR3XG6v2RyaVuPn3hjhq",
			extendedPublicKey:  "xpub6EsrSLgxCLDQSpRjCaRiojeys5c36xK69nHVUAt66CF6Ukxw6MpijJpkEupbwpHSte9EBwXsrx2mAGBqPyVK6dhxARffciPvVTodmUFCsy2",
		},
	}

	for _, testCase := range testCases {
		seed := seedFromMnemonic(testCase.mnemonic)
		key, err := bip32.NewMasterKey(seed)
		assert.Nil(t, err)

		radixKey, radixDeriveErr := deriveRadixPath(key)
		assert.Nil(t, radixDeriveErr)

		assert.Equal(t, testCase.extendedPrivateKey, radixKey.String(), testCase.name)
		assert.Equal(t, testCase.extendedPublicKey, radixKey.PublicKey().String(), testCase.name)
	}
}

func TestPublicKeysFromIndex(t *testing.T) {
	testCases := []struct {
		name              string
		mnemonic          string
		index             uint32
		extendedPublicKey string
	}{
		{
			name:              "Seed one",
			mnemonic:          "legal winner thank year wave sausage worth useful legal winner thank yellow",
			index:             0,
			extendedPublicKey: "036d39bd3894fa2193f1ffc62236bfadf3d3c051e8fe9ca5cc02677ea5e1ad34e8",
		},
		{
			name:              "Seed one",
			mnemonic:          "legal winner thank year wave sausage worth useful legal winner thank yellow",
			index:             3,
			extendedPublicKey: "028d31597419a690f369a079dfc54276b643836189a375b56d8c1983bffbb53c36",
		},
		{
			name:              "Seed two",
			mnemonic:          "special obey degree parrot lumber funny celery resist broken bundle robust popular",
			index:             1,
			extendedPublicKey: "0323e00ce137b54d7b3b0bafb997f3873f163c4d0fd508f064b5c5049b0a438e71",
		},
		{
			name:              "Seed two",
			mnemonic:          "special obey degree parrot lumber funny celery resist broken bundle robust popular",
			index:             2,
			extendedPublicKey: "033dee93777894e0f7829beaf07a7780257d10d1d048eec1498995c00f1a48f810",
		},
	}

	for _, testCase := range testCases {
		seed := seedFromMnemonic(testCase.mnemonic)
		key, err := bip32.NewMasterKey(seed)
		assert.Nil(t, err)

		radixKey, radixDeriveErr := deriveRadixPath(key)
		assert.Nil(t, radixDeriveErr)

		addressKey, err := getKeyForIndex(radixKey, testCase.index)
		assert.Nil(t, err)

		assert.Equal(t, testCase.extendedPublicKey, hex.EncodeToString(addressKey.PublicKey().Key), fmt.Sprintf("%s, Should be able to derive public key with index %d", testCase.name, testCase.index))
	}
}

func TestPublicAddressesFromIndex(t *testing.T) {
	testCases := []struct {
		name     string
		mnemonic string
		index    uint32
		address  string
		hrp      string
	}{
		{
			name:     "Seed one",
			mnemonic: "special obey degree parrot lumber funny celery resist broken bundle robust popular",
			index:    0,
			address:  "rdx1qspx2cw9jfsz3hxx0n2cl08llhttqh73c0pflacfdup308c3cyh87ecd94yfq",
			hrp:      "rdx",
		},
		{
			name:     "Seed one",
			mnemonic: "special obey degree parrot lumber funny celery resist broken bundle robust popular",
			index:    1,
			address:  "rdx1qspj8cqvuymm2ntm8v96lwvh7wrn793uf58a2z8svj6u2pympfpcuug0wfmm4",
			hrp:      "rdx",
		},
		{
			name:     "Seed two",
			mnemonic: "foam sketch royal kiwi field inject undo horn decorate usual gospel pottery",
			index:    0,
			address:  "rdx1qspc3lnyawt4xd66sedhk22a27uuxh6x4hfvl5plt62z7rzdstq7phsxu45a4",
			hrp:      "rdx",
		},
		{
			name:     "Seed two",
			mnemonic: "foam sketch royal kiwi field inject undo horn decorate usual gospel pottery",
			index:    1,
			address:  "rdx1qspr3d7rwmvn2r0kmf6adrzjv3p3qs02s4dqqnr5j5y59q0xnmehn8g5pk80p",
			hrp:      "rdx",
		},
		{
			name:     "Seed two",
			mnemonic: "foam sketch royal kiwi field inject undo horn decorate usual gospel pottery",
			index:    2,
			address:  "rdx1qspcj2rhh2z2mdrfnm7dw77caxhnrx3h55rdnf7x4kpaf7jj53a2degnv6kwz",
			hrp:      "rdx",
		},
		{
			name:     "Seed two",
			mnemonic: "foam sketch royal kiwi field inject undo horn decorate usual gospel pottery",
			index:    0,
			address:  "tdx1qspc3lnyawt4xd66sedhk22a27uuxh6x4hfvl5plt62z7rzdstq7phs8sqxdk",
			hrp:      "tdx",
		},
		{
			name:     "Seed two",
			mnemonic: "foam sketch royal kiwi field inject undo horn decorate usual gospel pottery",
			index:    1,
			address:  "tdx1qspr3d7rwmvn2r0kmf6adrzjv3p3qs02s4dqqnr5j5y59q0xnmehn8g4dr4lz",
			hrp:      "tdx",
		},
		{
			name:     "Seed two",
			mnemonic: "foam sketch royal kiwi field inject undo horn decorate usual gospel pottery",
			index:    2,
			address:  "tdx1qspcj2rhh2z2mdrfnm7dw77caxhnrx3h55rdnf7x4kpaf7jj53a2degjq0y7p",
			hrp:      "tdx",
		},
	}

	for _, testCase := range testCases {
		seed := seedFromMnemonic(testCase.mnemonic)
		key, err := bip32.NewMasterKey(seed)
		assert.Nil(t, err)

		radixKey, radixDeriveErr := deriveRadixPath(key)
		assert.Nil(t, radixDeriveErr)

		addressKey, err := getKeyForIndex(radixKey, testCase.index)
		assert.Nil(t, err)

		address, addressErr := deriveAddressFromPubKey(addressKey.PublicKey().Key, testCase.hrp)
		assert.Nil(t, addressErr)

		assert.Equal(t, testCase.address, address, fmt.Sprintf("%s, Should be able to derive public key with index %d", testCase.name, testCase.index))
	}
}
