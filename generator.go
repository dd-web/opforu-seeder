package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

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

var SlugAlphabet = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i",
	"j", "k", "l", "m", "n", "o", "p", "q", "r",
	"s", "t", "u", "v", "x", "y", "z", "-", "_",
	"0", "1", "2", "3", "4", "5", "6", "7", "8",
	"9",
}

// random int between min and max
func RandomIntBetween(min, max int) int {
	return min + rand.Intn(max-min)
}

// returns a random category of words
func SelectRandomCategory() *[]string {
	return categories[RandomIntBetween(0, len(categories)-1)]
}

// returns a random word from a category
func SelectRandomWord(c *[]string) string {
	str := (*c)[RandomIntBetween(0, len(*c)-1)]
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
	prefix := word_partials_prefix[RandomIntBetween(0, len(word_partials_prefix)-1)]
	suffix := word_partials_suffix[RandomIntBetween(0, len(word_partials_suffix)-1)]
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

	return word + "@" + words_domains[RandomIntBetween(0, len(words_domains)-1)] + "." + words_subdomains[RandomIntBetween(0, len(words_subdomains)-1)]
}

// return random lorem word
func GetLoremWord() string {
	return words_lorem[RandomIntBetween(0, len(words_lorem)-1)]
}

// return random sentence
func GetSentence() string {
	wc := RandomIntBetween(15, 20)
	sentence := ""
	for i := 0; i < wc; i++ {
		sentence = sentence + GetLoremWord() + " "
	}
	return sentence
}

// return random paragraph
func GetParagraph() string {
	sc := RandomIntBetween(3, 10)
	paragraph := "<p>"
	for i := 0; i < sc; i++ {
		paragraph = paragraph + GetSentence()
	}
	return paragraph + "</p>"
}

// return random number of paragraphs between min and max
func GetParagraphsBetween(min, max int) string {
	pc := RandomIntBetween(min, max)
	paragraphs := "<div class=\"thread-body\">"
	for i := 0; i < pc; i++ {
		paragraphs = paragraphs + GetParagraph()
	}
	return paragraphs + "</div>"
}

// weighted roles
func GetWeightedRole() string {
	num := RandomIntBetween(0, 100)
	if num < 85 {
		return "public"
	} else if num < 97 {
		return "mod"
	} else {
		return "admin"
	}
}

// weighted account status
func GetWeightedAccountStatus() string {
	num := RandomIntBetween(0, 100)
	if num < 92 {
		return "active"
	} else if num < 95 {
		return "suspended"
	} else if num < 98 {
		return "banned"
	} else {
		return "deleted"
	}
}

// weighted thread status
func GetWeightedThreadStatus() string {
	num := RandomIntBetween(0, 100)
	if num < 90 {
		return "open"
	} else if num < 95 {
		return "closed"
	} else if num < 97 {
		return "archived"
	} else {
		return "deleted"
	}
}

// weighted identity role
func GetWeightedIdentityRole() string {
	num := RandomIntBetween(0, 100)
	if num < 95 {
		return "public"
	} else {
		return "mod"
	}
}

// weighted identity status
func GetWeightedIdentityStatus() string {
	num := RandomIntBetween(0, 100)
	if num < 92 {
		return "active"
	} else if num < 98 {
		return "suspended"
	} else {
		return "banned"
	}
}

// return random slug between min and max
func GetSlug(min, max int) string {
	slugLen := RandomIntBetween(min, max)
	slugStr := ""
	for i := 0; i < slugLen; i++ {
		slugStr = slugStr + SlugAlphabet[RandomIntBetween(0, len(SlugAlphabet)-1)]
	}
	return slugStr
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

// return identity style string
func GetIdentityStyle() string {
	return "variant-" + aliasPrefixes[RandomIntBetween(0, len(aliasPrefixes)-1)] + "-" + aliasSuffixes[RandomIntBetween(0, len(aliasSuffixes)-1)]
}

// return between min and max tags
func GetRandomTags(min, max int) []string {
	tagCount := RandomIntBetween(min, max)
	tags := []string{}
	for i := 0; i < tagCount; i++ {
		tags = append(tags, GetLoremWord())
	}
	return tags
}

// pair of image url sections
func GetRandomImagePair() (string, string) {
	ix := RandomIntBetween(0, len(imageSourceSizes)-1)
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
	return imageFileExtensions[RandomIntBetween(0, len(imageFileExtensions)-1)]
}

// get random video file extension
func GetRandomVideoExt() string {
	return videoFileExtensions[RandomIntBetween(0, len(videoFileExtensions)-1)]
}

// get random media type
func GetRandomMediaType() string {
	num := RandomIntBetween(0, 100)
	if num > 85 {
		return "video"
	}
	return "image"
}

// get random media extension by type
func GetRandomExtByType(mediaType string) string {
	switch mediaType {
	case "image":
		return GetRandomImageExt()
	case "video":
		return GetRandomVideoExt()
	}

	return ""
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

// default password to use for dummy accounts
func GetDefaultPassword() string {
	return "123"
}
