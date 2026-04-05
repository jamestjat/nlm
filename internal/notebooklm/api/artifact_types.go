package api

// ArtifactTypeCode identifies the artifact type in CREATE_ARTIFACT (R7cb6c) RPC calls.
type ArtifactTypeCode int

const (
	ArtifactTypeAudio     ArtifactTypeCode = 1
	ArtifactTypeReport    ArtifactTypeCode = 2 // Briefing Doc, Study Guide, Blog Post, etc.
	ArtifactTypeVideo     ArtifactTypeCode = 3
	ArtifactTypeQuiz      ArtifactTypeCode = 4 // Also used for flashcards
	ArtifactTypeMindMap   ArtifactTypeCode = 5
	ArtifactTypeInfograph ArtifactTypeCode = 7
	ArtifactTypeSlideDeck ArtifactTypeCode = 8
	ArtifactTypeDataTable ArtifactTypeCode = 9
)

// ArtifactStatus represents the processing status of an artifact.
type ArtifactStatus int

const (
	ArtifactStatusProcessing ArtifactStatus = 1
	ArtifactStatusPending    ArtifactStatus = 2
	ArtifactStatusCompleted  ArtifactStatus = 3
	ArtifactStatusFailed     ArtifactStatus = 4
)

// AudioFormat selects the audio overview format.
type AudioFormat int

const (
	AudioFormatDeepDive  AudioFormat = 1
	AudioFormatBrief     AudioFormat = 2
	AudioFormatCritique  AudioFormat = 3
	AudioFormatDebate    AudioFormat = 4
)

// AudioLength selects the audio overview length.
type AudioLength int

const (
	AudioLengthShort   AudioLength = 1
	AudioLengthDefault AudioLength = 2
	AudioLengthLong    AudioLength = 3
)

// VideoFormat selects the video overview format.
type VideoFormat int

const (
	VideoFormatExplainer VideoFormat = 1
	VideoFormatBrief     VideoFormat = 2
	VideoFormatCinematic VideoFormat = 3
)

// VideoStyle selects the video visual style.
type VideoStyle int

const (
	VideoStyleAutoSelect  VideoStyle = 1
	VideoStyleCustom      VideoStyle = 2
	VideoStyleClassic     VideoStyle = 3
	VideoStyleWhiteboard  VideoStyle = 4
	VideoStyleKawaii      VideoStyle = 5
	VideoStyleAnime       VideoStyle = 6
	VideoStyleWatercolor  VideoStyle = 7
	VideoStyleRetroPrint  VideoStyle = 8
	VideoStyleHeritage    VideoStyle = 9
	VideoStylePaperCraft  VideoStyle = 10
)

// QuizQuantity controls number of quiz/flashcard items.
type QuizQuantity int

const (
	QuizQuantityFewer    QuizQuantity = 1
	QuizQuantityStandard QuizQuantity = 2
)

// QuizDifficulty controls quiz/flashcard difficulty.
type QuizDifficulty int

const (
	QuizDifficultyEasy   QuizDifficulty = 1
	QuizDifficultyMedium QuizDifficulty = 2
	QuizDifficultyHard   QuizDifficulty = 3
)

// InfographicOrientation controls infographic layout.
type InfographicOrientation int

const (
	InfographicOrientationLandscape InfographicOrientation = 1
	InfographicOrientationPortrait  InfographicOrientation = 2
	InfographicOrientationSquare    InfographicOrientation = 3
)

// InfographicDetail controls infographic verbosity.
type InfographicDetail int

const (
	InfographicDetailConcise  InfographicDetail = 1
	InfographicDetailStandard InfographicDetail = 2
	InfographicDetailDetailed InfographicDetail = 3
)

// InfographicStyle selects the infographic visual style.
type InfographicStyle int

const (
	InfographicStyleAutoSelect    InfographicStyle = 1
	InfographicStyleSketchNote    InfographicStyle = 2
	InfographicStyleProfessional  InfographicStyle = 3
	InfographicStyleBentoGrid     InfographicStyle = 4
	InfographicStyleEditorial     InfographicStyle = 5
	InfographicStyleInstructional InfographicStyle = 6
	InfographicStyleBricks        InfographicStyle = 7
	InfographicStyleClay          InfographicStyle = 8
	InfographicStyleAnime         InfographicStyle = 9
	InfographicStyleKawaii        InfographicStyle = 10
	InfographicStyleScientific    InfographicStyle = 11
)

// SlideDeckFormat selects the slide deck format.
type SlideDeckFormat int

const (
	SlideDeckFormatDetailed  SlideDeckFormat = 1
	SlideDeckFormatPresenter SlideDeckFormat = 2
)

// SlideDeckLength controls slide deck length.
type SlideDeckLength int

const (
	SlideDeckLengthDefault SlideDeckLength = 1
	SlideDeckLengthShort   SlideDeckLength = 2
)

// ExportType selects the export destination.
type ExportType int

const (
	ExportTypeDocs   ExportType = 1
	ExportTypeSheets ExportType = 2
)

// SourceStatus represents processing status of a source.
type SourceStatus int

const (
	SourceStatusProcessing SourceStatus = 1
	SourceStatusReady      SourceStatus = 2
	SourceStatusError      SourceStatus = 3
	SourceStatusPreparing  SourceStatus = 5
)

// ShareViewLevel controls what shared viewers can access.
type ShareViewLevel int

const (
	ShareViewLevelFullNotebook ShareViewLevel = 0
	ShareViewLevelChatOnly     ShareViewLevel = 1
)
