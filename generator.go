package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

var HrSplit string = "\n-----------------------------------------------------\n"

// @TODO maybe better to use byte arrays for performance - char codes via rune (int32)
// probably wont matter if it's lorem specifically or not, in which case the dataset can
// be significantly reduced while also increasing variance of the output
// lorem ipsum generator
var words_lorem = []string{"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"a", "ac", "accumsan", "ad", "aenean", "aliquam", "aliquet", "ante",
	"aptent", "arcu", "at", "auctor", "augue", "bibendum", "blandit",
	"class", "commodo", "condimentum", "congue", "consequat", "conubia",
	"convallis", "cras", "cubilia", "curabitur", "curae", "cursus",
	"dapibus", "diam", "dictum", "dictumst", "dignissim", "dis", "donec",
	"dui", "duis", "efficitur", "egestas", "eget", "eleifend", "elementum",
	"enim", "erat", "eros", "est", "et", "etiam", "eu", "euismod", "ex",
	"facilisi", "facilisis", "fames", "faucibus", "felis", "fermentum",
	"feugiat", "finibus", "fringilla", "fusce", "gravida", "habitant",
	"habitasse", "hac", "hendrerit", "himenaeos", "iaculis", "id",
	"imperdiet", "in", "inceptos", "integer", "interdum", "justo",
	"lacinia", "lacus", "laoreet", "lectus", "leo", "libero", "ligula",
	"litora", "lobortis", "luctus", "maecenas", "magna", "magnis",
	"malesuada", "massa", "mattis", "mauris", "maximus", "metus", "mi",
	"molestie", "mollis", "montes", "morbi", "mus", "nam", "nascetur",
	"natoque", "nec", "neque", "netus", "nibh", "nisi", "nisl", "non",
	"nostra", "nulla", "nullam", "nunc", "odio", "orci", "ornare",
	"parturient", "pellentesque", "penatibus", "per", "pharetra",
	"phasellus", "placerat", "platea", "porta", "porttitor", "posuere",
	"potenti", "praesent", "pretium", "primis", "proin", "pulvinar",
	"purus", "quam", "quis", "quisque", "rhoncus", "ridiculus", "risus",
	"rutrum", "sagittis", "sapien", "scelerisque", "sed", "sem", "semper",
	"senectus", "sociosqu", "sodales", "sollicitudin", "suscipit",
	"suspendisse", "taciti", "tellus", "tempor", "tempus", "tincidunt",
	"torquent", "tortor", "tristique", "turpis", "ullamcorper", "ultrices",
	"ultricies", "urna", "ut", "varius", "vehicula", "vel", "velit",
	"venenatis", "vestibulum", "vitae", "vivamus", "viverra", "volutpat",
	"vulputate"}

// color words
var words_colors = []string{
	"yellow",
	"red",
	"green",
	"blue",
	"orange",
	"purple",
	"pink",
	"black",
	"white",
	"brown",
	"gray",
	"silver",
	"gold",
	"bronze",
	"rainbow",
}

// sports words
var words_sports = []string{
	"football",
	"soccer",
	"basketball",
	"hockey",
	"baseball",
	"tennis",
	"golf",
	"rugby",
	"volleyball",
	"cricket",
	"badminton",
	"bowling",
	"boxing",
	"curling",
	"handball",
	"polo",
}

// animal words
var words_animals = []string{
	"dog",
	"cat",
	"bird",
	"fish",
	"horse",
	"cow",
	"pig",
	"sheep",
	"chicken",
	"duck",
	"goat",
	"turkey",
	"rabbit",
	"deer",
	"bear",
	"lion",
	"tiger",
	"elephant",
	"monkey",
}

// hobby words
var words_hobbies = []string{
	"fishing",
	"sailing",
	"swimming",
	"running",
	"jogging",
	"walking",
	"climbing",
	"skiing",
	"skating",
	"surfing",
	"skateboarding",
	"biking",
	"painting",
	"reading",
	"writing",
	"cooking",
}

// cities
var words_cities = []string{
	"new york",
	"los angeles",
	"chicago",
	"houston",
	"phoenix",
	"philadelphia",
	"san antonio",
	"san diego",
	"dallas",
	"san jose",
	"austin",
	"jacksonville",
	"fort worth",
	"san francisco",
}

// adjective words
var words_adjectives = []string{
	"happy",
	"sad",
	"angry",
	"mad",
	"glad",
	"funny",
	"cool",
	"hot",
	"cold",
	"gay",
	"fast",
	"slow",
	"quick",
	"smart",
	"stupid",
	"tall",
}

// noun words
var words_nouns = []string{
	"lover",
	"giver",
	"taker",
	"maker",
	"builder",
	"destroyer",
	"creator",
	"player",
	"fighter",
	"writer",
	"reader",
	"thinker",
	"doer",
	"worker",
	"helper",
	"smasher",
}

// email domains
var words_domains = []string{
	"gmail",
	"hotmail",
	"yahoo",
	"msn",
	"outlook",
	"live",
	"icloud",
	"protonmail",
	"zoho",
	"yandex",
}

// subdomains
var words_subdomains = []string{
	"com",
	"org",
	"net",
	"co",
	"edu",
	"gov",
	"mil",
	"int",
	"tv",
	"info",
	"biz",
}

// word prefixes
var word_partials_prefix = []string{
	"yu",
	"re",
	"gr",
	"bl",
	"or",
	"pu",
	"pi",
	"bl",
	"ph",
	"wh",
	"br",
	"gr",
	"si",
	"go",
	"za",
	"ra",
	"to",
	"dal",
	"dav",
	"wil",
	"jac",
	"for",
	"san",
	"los",
	"chi",
	"hou",
	"phi",
	"san",
	"tin",
	"xed",
	"fun",
	"coo",
	"hot",
	"col",
	"qui",
}

// word suffixes
var word_partials_suffix = []string{
	"ow",
	"ed",
	"en",
	"ue",
	"le",
	"er",
	"ty",
	"ck",
	"ly",
	"ne",
	"ny",
	"by",
	"ty",
	"ze",
	"ct",
	"ine",
	"ing",
	"ter",
	"der",
	"ker",
	"ler",
	"per",
	"yer",
	"mer",
	"ger",
	"ver",
	"ser",
	"zer",
	"ner",
	"ber",
	"der",
}

// word categories
var categories = []*[]string{
	&words_colors,
	&words_sports,
	&words_animals,
	&words_hobbies,
	&words_cities,
	&words_adjectives,
	&words_nouns,
}

// image base url
var imageBaseUrl = "https://picsum.photos/id/"

// some image sizes
var imageSourceSizes = []string{
	// 4:3
	"640/480",
	"800/600",
	"1024/768",
	// 16:9
	"640/360",
	"800/450",
	"960/540",
}

// image thumbnail sizes
var imageThumbnailSizes = []string{
	// 4:3
	"160/120",
	"200/150",
	"320/240",
	// 16:9
	"160/90",
	"320/180",
	"480/270",
}

// image file extensions
var imageFileExtensions = []string{
	"jpg",
	"jpeg",
	"png",
	"gif",
	"webp",
	"bmp",
	"svg",
}

// video flle extensions
var videoFileExtensions = []string{
	"mp4",
	"webm",
	"ogg",
	"avi",
	"mov",
}

var SlugAlphabet string = "abcdefghijklmnopqrstuvwxyz0123456789"

// random int between min and max
func RandomIntBetween(min, max int) int {
	return rand.Intn(max-min) + min
}

// returns a random category of words
func SelectRandomCategory() *[]string {
	return categories[RandomIntBetween(0, len(categories))]
}

// returns a random word from a category
func SelectRandomWord(c *[]string) string {
	str := (*c)[RandomIntBetween(0, len(*c))]
	return strings.ReplaceAll(str, " ", "")
}

// returns any random word from the categories of word - not field specific case words
func SelectAnyWord() string {
	return SelectRandomWord(SelectRandomCategory())
}

// return a random number of letters between min and max
func RandomLettersBetween(min, max int) string {
	strlen := RandomIntBetween(min, max)
	str := ""
	for i := 0; i < strlen; i++ {
		str = str + string(rune(RandomIntBetween(97, 122)))
	}
	return str
}

// some random username
func GetUsername() string {
	prefix := word_partials_prefix[RandomIntBetween(0, len(word_partials_prefix))]
	suffix := word_partials_suffix[RandomIntBetween(0, len(word_partials_suffix))]
	between := RandomLettersBetween(0, 5)
	word := ""

	c := RandomIntBetween(0, 8)

	switch c {
	case 0:
		word = prefix + between + suffix
	case 1:
		word = strconv.Itoa(RandomIntBetween(0, 99)) + prefix + between + suffix + strconv.Itoa(RandomIntBetween(0, 99))
	case 2:
		word = prefix + between + suffix + strconv.Itoa(RandomIntBetween(0, 99))
	case 3:
		word = prefix + between + strconv.Itoa(RandomIntBetween(0, 999)) + suffix
	case 4:
		word = prefix + between + strconv.Itoa(RandomIntBetween(0, 999)) + suffix + strconv.Itoa(RandomIntBetween(0, 99))
	case 5:
		word = prefix + strconv.Itoa(RandomIntBetween(0, 999)) + between + suffix + strconv.Itoa(RandomIntBetween(0, 999))
	case 6:
		word = prefix + strconv.Itoa(RandomIntBetween(0, 999)) + between + suffix
	case 7:
		word = strconv.Itoa(RandomIntBetween(0, 99)) + prefix + strconv.Itoa(RandomIntBetween(0, 999)) + between + suffix
	case 8:
		word = strconv.Itoa(RandomIntBetween(0, 99)) + prefix + strconv.Itoa(RandomIntBetween(0, 99)) + between + strconv.Itoa(RandomIntBetween(0, 99)) + suffix
	}

	return word
}

// some random email
func GetEmail() string {
	word := SelectAnyWord() + SelectAnyWord()
	c := RandomIntBetween(0, 10)

	switch c {
	case 0:
		word = word + strconv.Itoa(RandomIntBetween(0, 9))
	case 1:
		word = word + strconv.Itoa(RandomIntBetween(0, 99))
	case 2:
		word = word + strconv.Itoa(RandomIntBetween(0, 999))
	case 3:
		word = strconv.Itoa(RandomIntBetween(0, 999)) + word
	case 4:
		word = strconv.Itoa(RandomIntBetween(0, 99)) + word
	case 5:
		word = strconv.Itoa(RandomIntBetween(0, 9)) + word
	case 6:
		word = strconv.Itoa(RandomIntBetween(0, 9)) + word + strconv.Itoa(RandomIntBetween(0, 9))
	case 7:
		word = strconv.Itoa(RandomIntBetween(0, 99)) + word + strconv.Itoa(RandomIntBetween(0, 99))
	default:
		break
	}

	return word + "@" + words_domains[RandomIntBetween(0, len(words_domains))] + "." + words_subdomains[RandomIntBetween(0, len(words_subdomains))]
}

// return random lorem word
func GetLoremWord() string {
	return words_lorem[RandomIntBetween(0, len(words_lorem))]
}

// return random sentence
func GetSentence() string {
	wc := RandomIntBetween(6, 14)
	sentence := ""
	for i := 0; i < wc; i++ {
		sentence = sentence + GetLoremWord() + " "
	}
	return sentence
}

// return random paragraph
func GetParagraph() string {
	sc := RandomIntBetween(1, 6)
	paragraph := "<p>"
	for i := 0; i < sc; i++ {
		paragraph = paragraph + GetSentence()
	}
	return paragraph + "</p> "
}

// return random number of paragraphs between min and max
func GetParagraphsBetween(min, max int) string {
	pc := RandomIntBetween(min, max)
	paragraphs := "<div class=\"thread-body\"> "
	for i := 0; i < pc; i++ {
		paragraphs = paragraphs + GetParagraph()
	}
	return paragraphs + "</div> "
}

// weighted roles
func GetWeightedRole() AccountRole {
	num := RandomIntBetween(0, 100)
	if num < 90 {
		return AccountRoleUser
	} else if num < 93 {
		return AccountRoleMod
	} else if num < 98 {
		return AccountRoleAdmin
	} else {
		return AccountRolePublic
	}
}

// weighted account status
func GetWeightedAccountStatus() AccountStatus {
	num := RandomIntBetween(0, 100)
	if num < 92 {
		return AccountStatusActive
	} else if num < 95 {
		return AccountStatusSuspended
	} else if num < 98 {
		return AccountStatusBanned
	} else {
		return AccountStatusDeleted
	}
}

// weighted thread status
func GetWeightedThreadStatus() ThreadStatus {
	num := RandomIntBetween(0, 100)
	if num < 90 {
		return ThreadStatusOpen
	} else if num < 95 {
		return ThreadStatusClosed
	} else if num < 97 {
		return ThreadStatusArchived
	} else {
		return ThreadStatusDeleted
	}
}

// weighted identity status
func GetWeightedIdentityStatus() IdentityStatus {
	num := RandomIntBetween(0, 100)
	if num < 92 {
		return IdentityStatusActive
	} else if num < 98 {
		return IdentityStatusSuspended
	} else {
		return IdentityStatusBanned
	}
}

// weighted thread role (creator roles are specificly applied which is why they aren't here)
func GetWeightedThreadRole() ThreadRole {
	num := RandomIntBetween(0, 100)
	if num < 85 {
		return ThreadRoleUser
	} else {
		return ThreadRoleMod
	}
}

// return random slug between min and max
func GetSlug(min, max int) string {
	slugLen := RandomIntBetween(min, max)
	slug, _ := gonanoid.Generate(SlugAlphabet, slugLen)
	return slug
}

// identity alias prefixes
var aliasPrefixes = []string{
	"filled",
	"ghost",
	"soft",
	"glass",
}

// identity alias suffixes
var aliasSuffixes = []string{
	"primary",
	"secondary",
	"tertiary",
	"success",
	"warning",
	"error",
	"surface",
}

type ByteSize int

const (
	BYTE ByteSize = 8 << iota
	KB   ByteSize = 1 << (10 * iota)
	MB   ByteSize = 1 << (10 * iota)
	GB   ByteSize = 1 << (10 * iota)
	TB   ByteSize = 1 << (10 * iota)
)

var minFileSize int = int(BYTE) * 64 // 512 bytes
var maxFileSize int = int(MB) * 16   // 16 megabytes

var avatarFileSizeMultiplier int = 4

func GetRandomFileSize() int {
	return RandomIntBetween(minFileSize, maxFileSize)
}

// takes the number of total bytes and formats a string of the size
func FormatByteString(size int) string {
	if size < int(KB) {
		return strconv.Itoa(size) + "b"
	}

	if size < int(MB) {
		return strconv.FormatFloat(float64(size)/float64(KB), 'f', 2, 64) + "kb"
	}

	if size < int(GB) {
		return strconv.FormatFloat(float64(size)/float64(MB), 'f', 2, 64) + "mb"
	}

	if size < int(TB) {
		return strconv.FormatFloat(float64(size)/float64(GB), 'f', 2, 64) + "gb"
	}

	return strconv.FormatFloat(float64(size)/float64(TB), 'f', 2, 64) + "tb"
}

// return identity style string
func GetIdentityStyle() string {
	return "variant-" + aliasPrefixes[RandomIntBetween(0, len(aliasPrefixes))] + "-" + aliasSuffixes[RandomIntBetween(0, len(aliasSuffixes))]
}

// return between min and max tags
func GetRandomTags() []string {
	tagCount := RandomIntBetween(0, 6)
	words := map[string]int{}
	tags := []string{}

	for len(tags) < tagCount {
		w := GetLoremWord()
		if ok := words[w]; ok == 0 {
			words[w] = 1
			tags = append(tags, w)
		}
	}

	return tags
}

// md5 sample: AUqKpTEpkZ8OtOxqyXDPXw==
// sha256 sample: 4bt2e7eL3afWpEuEn1Yog14pea--PVGA4d06N0F7Tug=
// returns md5, sha256 base64 encoded checksums
func GetChecksumFromStr(str string) (string, string) {
	sumMD5, sumSHA256 := "", ""

	hashmd5 := md5.New()
	_, err := io.Copy(hashmd5, strings.NewReader(str))
	if err != nil {
		return "", ""
	}
	sumMD5 = base64.URLEncoding.EncodeToString(hashmd5.Sum(nil))

	hashsha256 := sha256.New()
	_, err = io.Copy(hashsha256, strings.NewReader(str))
	if err != nil {
		return "", ""
	}
	sumSHA256 = base64.URLEncoding.EncodeToString(hashsha256.Sum(nil))

	return sumMD5, sumSHA256
}

// generate and populate an asset source for deriving assets
func GenerateAssetSource(index int) *AssetSource {
	var width, height int = 0, 0     // dimensions
	var a_width, a_height int = 0, 0 // avatar dimnensions
	ix := RandomIntBetween(0, len(imageSourceSizes))

	kind := GetRandomAssetType()

	sizes := strings.Split(imageSourceSizes[ix], "/")
	width, _ = strconv.Atoi(sizes[0])
	height, _ = strconv.Atoi(sizes[1])

	a_size := strings.Split(imageThumbnailSizes[ix], "/")
	a_width, _ = strconv.Atoi(a_size[0])
	a_height, _ = strconv.Atoi(a_size[1])

	ts := time.Now().UTC()
	tsn := ts.UnixNano()

	sourceURL, avatarURL := FormatImageUrls(index)
	ext := GetRandomAssetExt(kind)

	rawSize := GetRandomFileSize()
	a_rawSize := rawSize / avatarFileSizeMultiplier

	cs_md5, cs_sha256 := GetChecksumFromStr(sourceURL)
	acs_md5, acs_sha256 := GetChecksumFromStr(avatarURL)

	sourceFileCtx := &FileCtx{
		ServerFileName: fmt.Sprintf("%d", tsn),
		Height:         uint16(height),
		Width:          uint16(width),
		FileSize:       uint32(rawSize),
		URL:            sourceURL,
		Extension:      ext,
		HashMD5:        cs_md5,
		HashSHA256:     cs_sha256,
	}

	avatarFileCtx := &FileCtx{
		ServerFileName: fmt.Sprintf("a-%d", tsn),
		Height:         uint16(a_height),
		Width:          uint16(a_width),
		FileSize:       uint32(a_rawSize),
		URL:            avatarURL,
		Extension:      ext,
		HashMD5:        acs_md5,
		HashSHA256:     acs_sha256,
	}

	details := &AssetSourceDetails{
		Source: sourceFileCtx,
		Avatar: avatarFileCtx,
	}

	src := &AssetSource{
		ID:        primitive.NewObjectID(),
		Details:   details,
		AssetType: kind,
		Uploaders: []primitive.ObjectID{},
		CreatedAt: &ts,
		UpdatedAt: &ts,
	}

	return src
}

// pair of image url sections
func GetRandomImagePair() (string, string) {
	ix := RandomIntBetween(0, len(imageSourceSizes))
	return imageSourceSizes[ix], imageThumbnailSizes[ix]
}

// format url with id and size
func GetImageUrl(id int, size string) string {
	return imageBaseUrl + fmt.Sprintf("%d/%s", id, size)
}

// gets a pair of image urls from a random id
func FormatImageUrls(id int) (string, string) {
	source, thumb := GetRandomImagePair()
	return GetImageUrl(id, source), GetImageUrl(id, thumb)
}

// get random image file extension
func GetRandomImageExt() string {
	return imageFileExtensions[RandomIntBetween(0, len(imageFileExtensions))]
}

// get random video file extension
func GetRandomVideoExt() string {
	return videoFileExtensions[RandomIntBetween(0, len(videoFileExtensions))]
}

// get random AssetType (not string)
func GetRandomAssetType() AssetType {
	num := RandomIntBetween(0, 100)
	if num > 90 {
		return AssetTypeVideo
	}
	return AssetTypeImage
}

// gets a random asset extension based on type
func GetRandomAssetExt(at AssetType) string {
	switch at {
	case AssetTypeImage:
		return GetRandomImageExt()
	case AssetTypeVideo:
		return GetRandomVideoExt()
	default:
		return ""
	}
}

// md5 checksum of file
func GetFileChecksumMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// hashes a password with bcrypt
func HashPassword(plaintext string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// compare a hashed password with plaintext
func ComparePass(hashed, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
