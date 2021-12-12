package mp4

type FourCC [4]byte

type BoxType FourCC

var (
	AvcCBoxType = BoxType{'a', 'v', 'c', 'C'}
	AvcEBoxType = BoxType{'a', 'v', 'c', 'E'}
	BtrtBoxType = BoxType{'b', 't', 'r', 't'}
	ClapBoxType = BoxType{'c', 'l', 'a', 'p'}
	ColrBoxType = BoxType{'c', 'o', 'l', 'r'}
	CttsBoxType = BoxType{'c', 't', 't', 's'}
	DinfBoxType = BoxType{'d', 'i', 'n', 'f'}
	DrefBoxType = BoxType{'d', 'r', 'e', 'f'}
	DvcCBoxType = BoxType{'d', 'v', 'c', 'C'}
	DvvCBoxType = BoxType{'d', 'v', 'v', 'C'}
	DvwCBoxType = BoxType{'d', 'v', 'w', 'C'}
	ElngBoxType = BoxType{'e', 'l', 'n', 'g'}
	EncaBoxType = BoxType{'e', 'n', 'c', 'a'}
	EncsBoxType = BoxType{'e', 'n', 'c', 's'}
	EnctBoxType = BoxType{'e', 'n', 'c', 't'}
	EncvBoxType = BoxType{'e', 'n', 'c', 'v'}
	FreeBoxType = BoxType{'f', 'r', 'e', 'e'}
	FrmaBoxType = BoxType{'f', 'r', 'm', 'a'}
	FtypBoxType = BoxType{'f', 't', 'y', 'p'}
	HdlrBoxType = BoxType{'h', 'd', 'l', 'r'}
	HvcCBoxType = BoxType{'h', 'v', 'c', 'C'}
	HvcEBoxType = BoxType{'h', 'v', 'c', 'E'}
	MdatBoxType = BoxType{'m', 'd', 'a', 't'}
	MdhdBoxType = BoxType{'m', 'd', 'h', 'd'}
	MdiaBoxType = BoxType{'m', 'd', 'i', 'a'}
	MfhdBoxType = BoxType{'m', 'f', 'h', 'd'}
	MinfBoxType = BoxType{'m', 'i', 'n', 'f'}
	MoofBoxType = BoxType{'m', 'o', 'o', 'f'}
	MoovBoxType = BoxType{'m', 'o', 'o', 'v'}
	MvexBoxType = BoxType{'m', 'v', 'e', 'x'}
	MvhdBoxType = BoxType{'m', 'v', 'h', 'd'}
	NmhdBoxType = BoxType{'n', 'm', 'h', 'd'}
	PaspBoxType = BoxType{'p', 'a', 's', 'p'}
	PsshBoxType = BoxType{'p', 's', 's', 'h'}
	SaioBoxType = BoxType{'s', 'a', 'i', 'o'}
	SaizBoxType = BoxType{'s', 'a', 'i', 'z'}
	SchiBoxType = BoxType{'s', 'c', 'h', 'i'}
	SchmBoxType = BoxType{'s', 'c', 'h', 'm'}
	SencBoxType = BoxType{'s', 'e', 'n', 'c'}
	SinfBoxType = BoxType{'s', 'i', 'n', 'f'}
	SmhdBoxType = BoxType{'s', 'm', 'h', 'd'}
	StblBoxType = BoxType{'s', 't', 'b', 'l'}
	StcoBoxType = BoxType{'s', 't', 'c', 'o'}
	StdpBoxType = BoxType{'s', 't', 'd', 'p'}
	StscBoxType = BoxType{'s', 't', 's', 'c'}
	StsdBoxType = BoxType{'s', 't', 's', 'd'}
	StssBoxType = BoxType{'s', 't', 's', 's'}
	StszBoxType = BoxType{'s', 't', 's', 'z'}
	SttsBoxType = BoxType{'s', 't', 't', 's'}
	TencBoxType = BoxType{'t', 'e', 'n', 'c'}
	TfhdBoxType = BoxType{'t', 'f', 'h', 'd'}
	TkhdBoxType = BoxType{'t', 'k', 'h', 'd'}
	TrakBoxType = BoxType{'t', 'r', 'a', 'k'}
	TrafBoxType = BoxType{'t', 'r', 'a', 'f'}
	TrexBoxType = BoxType{'t', 'r', 'e', 'x'}
	TrunBoxType = BoxType{'t', 'r', 'u', 'n'}
	UuidBoxType = BoxType{'u', 'u', 'i', 'd'}
	UrlBoxType  = BoxType{'u', 'r', 'l', ' '}
	UrnBoxType  = BoxType{'u', 'r', 'n', ' '}
	VmhdBoxType = BoxType{'v', 'm', 'h', 'd'}

	DvavBoxType = BoxType{'d', 'v', 'a', 'v'}
	Dva1BoxType = BoxType{'d', 'v', 'a', '1'}
	DvheBoxType = BoxType{'d', 'v', 'h', 'e'}
	Dvh1BoxType = BoxType{'d', 'v', 'h', '1'}
	Avc1BoxType = BoxType{'a', 'v', 'c', '1'}
	Avc2BoxType = BoxType{'a', 'v', 'c', '2'}
	Avc3BoxType = BoxType{'a', 'v', 'c', '3'}
	Avc4BoxType = BoxType{'a', 'v', 'c', '4'}
	Hev1BoxType = BoxType{'h', 'e', 'v', '1'}
	Hvc1BoxType = BoxType{'h', 'v', 'c', '1'}

	Avc1FourCC = FourCC{'a', 'v', 'c', '1'}
	Avc2FourCC = FourCC{'a', 'v', 'c', '2'}
	Avc3FourCC = FourCC{'a', 'v', 'c', '3'}
	Avc4FourCC = FourCC{'a', 'v', 'c', '4'}
	CencFourCC = FourCC{'c', 'e', 'n', 'c'}
	DashFourCC = FourCC{'d', 'a', 's', 'h'}
	DvavFourCC = FourCC{'d', 'v', 'a', 'v'}
	Dva1FourCC = FourCC{'d', 'v', 'a', '1'}
	DvheFourCC = FourCC{'d', 'v', 'h', 'e'}
	Dvh1FourCC = FourCC{'d', 'v', 'h', '1'}
	Hev1FourCC = FourCC{'h', 'e', 'v', '1'}
	HintFourCC = FourCC{'h', 'i', 'n', 't'}
	Hvc1FourCC = FourCC{'h', 'v', 'c', '1'}
	Iso2FourCC = FourCC{'i', 's', 'o', '2'}
	Iso6FourCC = FourCC{'i', 's', 'o', '6'}
	IsomFourCC = FourCC{'i', 's', 'o', 'm'}
	MetaFourCC = FourCC{'m', 'e', 't', 'a'}
	MsdhFourCC = FourCC{'m', 's', 'd', 'h'}
	SounFourCC = FourCC{'s', 'o', 'u', 'n'}
	VideFourCC = FourCC{'v', 'i', 'd', 'e'}

	NclcFourCC = FourCC{'n', 'c', 'l', 'c'}
	NclxFourCC = FourCC{'n', 'c', 'l', 'x'}
	RiccFourCC = FourCC{'r', 'I', 'C', 'C'}
	ProfFourCC = FourCC{'p', 'r', 'o', 'f'}

	SampleEncryptionBoxUserType = UserType{0xA2, 0x39, 0x4F, 0x52, 0x5A, 0x9B, 0x4F, 0x14, 0xA2, 0x44, 0x6C, 0x42, 0x7C, 0x64, 0x8D, 0xF4}
)
