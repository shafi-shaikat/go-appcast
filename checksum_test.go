package appcast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChecksum(t *testing.T) {
	c := NewChecksum(Sha256, "test")
	assert.IsType(t, Checksum{}, *c)
	assert.Equal(t, Sha256, c.Algorithm)
	assert.Equal(t, "test", c.Source)
	assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", c.Result)
}

func TestGenerate(t *testing.T) {
	testCases := map[string][]string{
		"sourceforge_default.xml": {
			"47e507823eb16c4cf577020226f1e5567a5d83173ee15d3a1a4fbae737550211",
			"e54afe17c3256480b88d57f4d5af88cd868c19e983fcefef3e05bedb69ace5ff",
			"f7273db54a49878a933af479d7efc837",
		},
		"sourceforge_empty.xml": {
			"f3fd86c2fd623cd4439bcf41c00eab3de8815e7aa59ea7c635775e0fe6293159",
			"f1d5fda5146d51438658a21b39d79f83de1689e8fe7ad9494946d3704146b452",
			"47a8f9860e6ce31f4525cf29c62b5262",
		},
		"sourceforge_single.xml": {
			"a5c0d4b0ba7867267acd4c8e3e0beefd536095bf34188246a06795e4bf456b22",
			"1c4450daa13e5aa52b5a65e14bf71cefe7928f03141cab7c8aade1f939dbdda3",
			"46c0ce057d804cb76cdfb204336135ea",
		},
		"sparkle_attributes_as_elements.xml": {
			"8c42d7835109ff61fe85bba66a44689773e73e0d773feba699bceecefaf09359",
			"15e08d20c984c6462632401405d4c74651f8bbb6d8924dab29d57f21cd23fbda",
			"90444fe711048735501877fd54dbcbd3",
		},
		"sparkle_default_asc.xml": {
			"9f94a728eab952284b47cc52acfbbb64de71f3d38e5b643d1f3523ef84495d9f",
			"f5f9b5d1d55ea8e5260b7537e9a5ad7b8dc7d43610a184b4a063416a7ee88c40",
			"0247ff43c3df1a0c6c3f2bedf5f4be05",
		},
		"sparkle_default.xml": {
			"83c1fd76a250dd50334db793a0db5da7575fc83d292c7c58fd9d31d5bcef6566",
			"87007d361728a5f02452552a8245e7f918521d2fc8a28c039972616aa7abfadc",
			"56157a2dc1cec9dc02448223e31854fa",
		},
		"sparkle_incorrect_namespace.xml": {
			"2e66ef346c49a8472bf8bf26e6e778c5b4d494723223c84c35d9f272a7792430",
			"52c66bf81606819d16d69202ff6836d18e0d2fa9d817097f7bd57e7c8a5b6215",
			"da82e1a170325e28e4fc1ed94bacaa88",
		},
		"sparkle_invalid_pubdate.xml": {
			"e0273ccbce5a6fb6a5fe31b5edffb8173d88afa308566cf9b4373f3fed909705",
			"87007d361728a5f02452552a8245e7f918521d2fc8a28c039972616aa7abfadc",
			"d98e602b718c7949a88fd41d9cc28cc8",
		},
		"sparkle_invalid_version.xml": {
			"12c7827fed4cccb5c4bc77052d2c95b03c0e4943aa49c90f9f2e98bb8ab9b799",
			"6eb256a3d3c226146d985a712325bd488fba6dcd47ca1ea5a48bc5535bff5fc9",
			"6a2a0417379a4f70272165fe053c76d0",
		},
		"sparkle_multiple_enclosure.xml": {
			"7f62916d4d80cc9a784ffa1d2211488104c4578cc2704baaff48a96b4df00961",
			"0c927f077bcd492fcd574bf6689d2131d96d368c5d7516f6d4c6cc645e12114d",
			"27737897524ca35a512c0ef4d9cff44a",
		},
		"sparkle_no_releases.xml": {
			"befd99d96be280ca7226c58ef1400309905ad20d2723e69e829cf050e802afcf",
			"befd99d96be280ca7226c58ef1400309905ad20d2723e69e829cf050e802afcf",
			"6b2f1f5e0cea6005e5410c1d76cab0a3",
		},
		"sparkle_single.xml": {
			"ac649bebe55f84d85767072e3a1122778a04e03f56b78226bd57ab50ce9f9306",
			"b36b6e57c9aec8ffa913dbfb1b4dba01d2f9246c7eaa4cfd0084a574216cf4a0",
			"aa5b165d930d81645b5c14c66bf67957",
		},
		"sparkle_with_comments.xml": {
			"283ea10e6f7f81466beb85e055940765f308dfdd7fd3ee717a65a4e19b31b460",
			"68d1e75702d38e1d84807413366ea9c1bbcca614068965c0a63cf84e75dd9848",
			"159973849b349fa9b37d2287af8dd528",
		},
		"sparkle_without_namespaces.xml": {
			"ee2d28f74e7d557bd7259c0f24a261658a9f27a710308a5c539ab761dae487c1",
			"b5b2055b16c135670a885ccb8018705b5500bff1f5ed65ad79ec5903da47beec",
			"d4c80271cfff4ab0afc15f7699c2e376",
		},
		"unknown.xml": {
			"c29665078d79a8e67b37b46a51f2a34c6092719833ccddfdda6109fd8f28043c",
			"c29665078d79a8e67b37b46a51f2a34c6092719833ccddfdda6109fd8f28043c",
			"2340f9a888f7305f4636d4f70d3471b1",
		},
	}

	for filename, checkpoints := range testCases {
		content := string(getTestdata(filename))

		// SHA256
		c := &Checksum{Sha256, content, ""}
		assert.Equal(t, checkpoints[0], c.Generate(), fmt.Sprintf("Checksum doesn't match (Sha256): %s", filename))

		// SHA256 (Homebrew-Cask)
		c = &Checksum{Sha256HomebrewCask, content, ""}
		assert.Equal(t, checkpoints[1], c.Generate(), fmt.Sprintf("Checksum doesn't match (Sha256HomebrewCask): %s", filename))

		// MD5
		c = &Checksum{Md5, content, ""}
		assert.Equal(t, checkpoints[2], c.Generate(), fmt.Sprintf("Checksum doesn't match (Md5): %s", filename))
	}
}