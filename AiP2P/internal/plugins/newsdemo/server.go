package newsplugin

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"aip2p.org/internal/apphost"
)

//go:embed web/templates/*.html web/static/*
var webFS embed.FS

type App struct {
	storeRoot  string
	project    string
	version    string
	archive    string
	rulesPath  string
	writerPath string
	netPath    string
	listenAddr string
	templates  *template.Template
	staticFS   fs.FS
	loadIndex  func(storeRoot, project string) (Index, error)
	syncIndex  func(index *Index, archiveRoot string) error
	loadRules  func(path string) (SubscriptionRules, error)
	loadWriter func(path string) (WriterPolicy, error)
	loadNet    func(path string) (NetworkBootstrapConfig, error)
	loadSync   func(storeRoot string) (SyncRuntimeStatus, error)
	loadSuper  func(path string) (SyncSupervisorState, error)
	fetchLANBT func(ctx context.Context, value, expectedNetworkID string) (NetworkBootstrapResponse, error)
	options    AppOptions
}

type AppOptions struct {
	ContentRoutes      bool
	ContentAPIRoutes   bool
	ArchiveRoutes      bool
	HistoryAPIRoutes   bool
	NetworkRoutes      bool
	NetworkAPIRoutes   bool
	WriterPolicyRoutes bool
}

type NavItem struct {
	Name   string
	URL    string
	Active bool
}

type MetadataField struct {
	Label string
	Value string
}

type ModuleFocusCard struct {
	Title       string
	Description string
	Items       []string
}

type ModuleBoardColumn struct {
	Title       string
	Description string
	Path        string
	Count       int
	Posts       []Post
}

type ModuleReplyLane struct {
	Key         string
	Title       string
	Description string
	Count       int
	Replies     []Reply
}

type ModuleReplyHighlight struct {
	Key         string
	Title       string
	Description string
	Reply       *Reply
}

type CoordSection struct {
	Title       string
	Description string
	Items       []MetadataField
}

type CoordAction struct {
	Label string
	URL   string
}

type WorkbenchLane struct {
	Title       string
	Kind        string
	Description string
	Path        string
	EmptyCopy   string
	Posts       []Post
}

type WorkbenchCluster struct {
	Task            Post
	Related         []CoordRelatedGroup
	RegistryActions []CoordAction
	JumpPanel       *WorkspaceJumpPanel
}

type AgentWorkbenchCard struct {
	Agent Post
	Panel *CoordWorkbenchPanel
}

type RegistryWorkbenchCard struct {
	Kind          string
	Name          string
	URL           string
	ExternalURL   string
	Summary       []SummaryStat
	ModuleActions []CoordAction
	JumpPanel     *WorkspaceJumpPanel
	PrimaryCount  int
	ReplyCount    int
	ReactionCount int
}

type WorkstreamRelationPath struct {
	From          Post
	RelationKey   string
	RelationTitle string
	Kind          string
	Path          string
	Target        Post
	Reason        string
	Score         int
}

type WorkstreamRelationPanel struct {
	Title           string
	Description     string
	PrimaryPaths    []WorkstreamRelationPath
	SupportingPaths []WorkstreamRelationPath
}

type WorkspaceJumpPanel struct {
	Eyebrow           string
	Title             string
	Description       string
	PrimaryActions    []CoordAction
	SupportingActions []CoordAction
}

type HomePriorityCard struct {
	Eyebrow     string
	Title       string
	Description string
	URL         string
	Stats       []SummaryStat
	Actions     []CoordAction
}

type CoordRelatedGroup struct {
	Title string
	Kind  string
	Path  string
	Posts []Post
}

type CoordRelationGroup struct {
	Key            string
	Title          string
	Description    string
	DominantReason string
	Evidence       []string
	Kind           string
	Path           string
	Priority       int
	Score          int
	Items          []CoordRelationItem
}

type CoordRelationItem struct {
	Post     Post
	Evidence []string
	Score    int
}

type CoordWorkbenchPanel struct {
	Eyebrow     string
	Title       string
	Description string
	Groups      []CoordRelatedGroup
}

type CoordDetailFocusPanel struct {
	Eyebrow     string
	Title       string
	Description string
	Checklist   []string
	Actions     []CoordAction
}

type DirectoryItem struct {
	Name          string
	URL           string
	ExternalURL   string
	AssetCount    int
	ReplyCount    int
	ReactionCount int
	AvgTruth      string
}

type NodeStatusEntry struct {
	Label  string
	Value  string
	Detail string
	Tone   string
}

type NodeStatusCard struct {
	Label  string
	Value  string
	Detail string
	Tone   string
}

type NodeStatus struct {
	Summary       string
	SummaryTone   string
	SummaryDetail string
	NetworkStatus string
	NetworkTone   string
	NetworkDetail string
	Entries       []NodeStatusEntry
	Dashboard     []NodeStatusCard
}

type HomePageData struct {
	Project           string
	Version           string
	Posts             []Post
	LegacyPosts       []Post
	Now               time.Time
	ListenAddr        string
	AgentView         bool
	ShowNetworkWarn   bool
	Options           FeedOptions
	PageNav           []NavItem
	TopicFacets       []FeedFacet
	SourceFacets      []FeedFacet
	SortOptions       []SortOption
	WindowOptions     []TimeWindowOption
	PageSizeOptions   []PageSizeOption
	ActiveFilters     []ActiveFilter
	CoordGroupFacets  []FeedFacet
	SummaryStats      []SummaryStat
	PriorityCards     []HomePriorityCard
	WorkspaceStats    []SummaryStat
	WorkspaceActions  []CoordAction
	WorkspaceLanes    []WorkbenchLane
	WorkspaceClusters []WorkbenchCluster
	AgentWorkspace    []AgentWorkbenchCard
	TopicWorkspace    []RegistryWorkbenchCard
	SourceWorkspace   []RegistryWorkbenchCard
	TotalPostCount    int
	Pagination        PaginationState
	Subscriptions     SubscriptionRules
	NodeStatus        NodeStatus
}

type CollectionPageData struct {
	Project         string
	Version         string
	Kind            string
	Name            string
	Path            string
	DirectoryURL    string
	APIPath         string
	Now             time.Time
	Posts           []Post
	Options         FeedOptions
	PageNav         []NavItem
	SortOptions     []SortOption
	WindowOptions   []TimeWindowOption
	PageSizeOptions []PageSizeOption
	SideLabel       string
	SideFacets      []FeedFacet
	ActiveFilters   []ActiveFilter
	SummaryStats    []SummaryStat
	WorkspaceStats  []SummaryStat
	WorkspaceLanes  []WorkbenchLane
	JumpPanel       *WorkspaceJumpPanel
	RelationPanel   *WorkstreamRelationPanel
	ModuleActions   []CoordAction
	TotalPostCount  int
	Pagination      PaginationState
	ExternalURL     string
	NodeStatus      NodeStatus
}

type TypedCollectionPageData struct {
	Project          string
	Version          string
	AgentView        bool
	Title            string
	CoordType        string
	MetadataSchema   *CoordMetadataTypeSpec
	PublishGuide     *CoordPublishGuide
	Description      string
	ModuleBoard      []ModuleBoardColumn
	Path             string
	APIPath          string
	Now              time.Time
	Posts            []Post
	Options          FeedOptions
	PageNav          []NavItem
	SortOptions      []SortOption
	WindowOptions    []TimeWindowOption
	PageSizeOptions  []PageSizeOption
	TopicFacetLabel  string
	TopicFacets      []FeedFacet
	SourceFacetLabel string
	SourceFacets     []FeedFacet
	MetaFacetLabel   string
	MetaFacets       []FeedFacet
	ActiveFilters    []ActiveFilter
	SummaryStats     []SummaryStat
	TotalPostCount   int
	Pagination       PaginationState
	EmptyTitle       string
	EmptyCopy        string
	NodeStatus       NodeStatus
}

type DirectoryPageData struct {
	Project      string
	Version      string
	Kind         string
	Path         string
	APIPath      string
	Now          time.Time
	PageNav      []NavItem
	Items        []DirectoryItem
	SummaryStats []SummaryStat
	NodeStatus   NodeStatus
}

type PostPageData struct {
	Project               string
	Version               string
	AgentView             bool
	PageNav               []NavItem
	Post                  Post
	CoordType             string
	MetadataSchema        *CoordMetadataTypeSpec
	PublishGuide          *CoordPublishGuide
	CoordFields           []MetadataField
	CoordSections         []CoordSection
	JumpPanel             *WorkspaceJumpPanel
	DetailFocus           *CoordDetailFocusPanel
	CoordRelations        []CoordRelationGroup
	CoordWorkspace        *CoordWorkbenchPanel
	CoordRelated          []CoordRelatedGroup
	ModuleReplyLanes      []ModuleReplyLane
	ModuleReplyHighlights []ModuleReplyHighlight
	Replies               []Reply
	Reactions             []Reaction
	Related               []Post
	NodeStatus            NodeStatus
}

type ArchiveIndexPageData struct {
	Project       string
	Version       string
	PageNav       []NavItem
	Now           time.Time
	Days          []ArchiveDay
	SummaryStats  []SummaryStat
	Subscriptions SubscriptionRules
	NodeStatus    NodeStatus
}

type ArchiveDayPageData struct {
	Project       string
	Version       string
	PageNav       []NavItem
	Now           time.Time
	Day           string
	Days          []ArchiveDay
	Entries       []ArchiveEntry
	SummaryStats  []SummaryStat
	Subscriptions SubscriptionRules
	NodeStatus    NodeStatus
}

type ArchiveMessagePageData struct {
	Project    string
	Version    string
	PageNav    []NavItem
	Now        time.Time
	Entry      ArchiveEntry
	Content    string
	Thread     string
	RawURL     string
	DayURL     string
	Archive    string
	NodeStatus NodeStatus
}

type NetworkPageData struct {
	Project       string
	Version       string
	ListenAddr    string
	PageNav       []NavItem
	Now           time.Time
	NodeStatus    NodeStatus
	SyncStatus    SyncRuntimeStatus
	Supervisor    SyncSupervisorState
	LANPeers      []string
	LANBTAnchors  []LANBTAnchorStatus
	LANBTOverall  string
	LANBTHasMatch bool
}

type LANBTAnchorStatus struct {
	Peer        string
	Nodes       []string
	MatchedNode string
	Error       string
}

type NetworkBootstrapResponse struct {
	Project         string   `json:"project"`
	Version         string   `json:"version"`
	NetworkID       string   `json:"network_id"`
	PeerID          string   `json:"peer_id"`
	ListenAddrs     []string `json:"listen_addrs"`
	DialAddrs       []string `json:"dial_addrs"`
	BitTorrentNodes []string `json:"bittorrent_nodes,omitempty"`
}

func NewWithThemeAndOptions(storeRoot, project, version, archiveRoot, rulesPath, writerPath, netPath string, theme apphost.WebTheme, options AppOptions) (*App, error) {
	return newApp(storeRoot, project, version, archiveRoot, rulesPath, writerPath, netPath, theme, options)
}

func newApp(storeRoot, project, version, archiveRoot, rulesPath, writerPath, netPath string, theme apphost.WebTheme, options AppOptions) (*App, error) {
	if err := ensureRuntimeLayout(storeRoot, archiveRoot, rulesPath, writerPath, netPath); err != nil {
		return nil, err
	}
	funcs := template.FuncMap{
		"formatTime": func(t time.Time) string { return t.Format("2006-01-02 15:04 MST") },
		"formatOptionalTime": func(t *time.Time) string {
			if t == nil {
				return "none yet"
			}
			return t.Format("2006-01-02 15:04 MST")
		},
		"formatScore": func(value *float64) string {
			if value == nil {
				return "-"
			}
			return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(*value, 'f', 2, 64), "0"), ".")
		},
		"compactIdentity": compactIdentity,
		"isPublicKeyish":  isPublicKeyish,
		"displayArchivePath": func(value string) string {
			value = filepath.ToSlash(strings.TrimSpace(value))
			if value == "" {
				return ""
			}
			if idx := strings.Index(value, "/archive/"); idx >= 0 {
				return strings.TrimPrefix(value[idx+1:], "/")
			}
			if strings.HasPrefix(value, "archive/") {
				return value
			}
			return filepath.Base(value)
		},
		"join":                        strings.Join,
		"lower":                       strings.ToLower,
		"reactionLabel":               reactionLabel,
		"sourcePath":                  SourcePath,
		"topicPath":                   TopicPath,
		"coordPath":                   CoordPath,
		"coordAPIPath":                CoordAPIPath,
		"articlePreview":              ArticlePreviewForPost,
		"articleHasMore":              ArticleHasMoreForPost,
		"articleClass":                StandardArticleClassForPost,
		"standardArticleClasses":      StandardArticleClasses,
		"articleBody":                 ArticleBodyForPost,
		"articleIsCode":               ArticleIsCodeForPost,
		"articleCode":                 ArticleCodeForPost,
		"articleCodeLanguage":         ArticleCodeLanguageForPost,
		"articleIsJSON":               ArticleIsJSONForPost,
		"articleLinks":                ArticleLinksForPost,
		"publisherBio":                PublisherBioForPost,
		"replyPreview":                ReplyPreview,
		"replyHasMore":                ReplyHasMore,
		"moduleCollectionTitle":       ModuleCollectionTitle,
		"moduleCollectionDescription": ModuleCollectionDescription,
		"moduleCollectionChecklist":   ModuleCollectionChecklist,
		"moduleCollectionFocusCards":  ModuleCollectionFocusCards,
		"moduleCollectionQuickStats":  ModuleCollectionQuickStats,
		"modulePrimaryFields":         ModulePrimaryFieldsForPost,
		"moduleDetailQuickStats":      ModuleDetailQuickStats,
		"moduleDetailFocusCards":      ModuleDetailFocusCards,
	}
	tmpl, staticFS, err := loadThemeAssets(theme, funcs)
	if err != nil {
		return nil, err
	}
	return &App{
		storeRoot:  storeRoot,
		project:    project,
		version:    strings.TrimSpace(version),
		archive:    archiveRoot,
		rulesPath:  rulesPath,
		writerPath: writerPath,
		netPath:    netPath,
		templates:  tmpl,
		staticFS:   staticFS,
		loadIndex:  LoadIndex,
		syncIndex:  SyncMarkdownArchive,
		loadRules:  LoadSubscriptionRules,
		loadWriter: LoadWriterPolicy,
		loadNet:    LoadNetworkBootstrapConfig,
		loadSync:   loadSyncRuntimeStatus,
		loadSuper:  loadSyncSupervisorState,
		fetchLANBT: fetchNetworkBootstrapResponse,
		options:    options,
	}, nil
}

func (a *App) index() (Index, error) {
	index, err := a.loadIndex(a.storeRoot, a.project)
	if err != nil {
		return Index{}, err
	}
	if a.loadRules != nil {
		rules, err := a.loadRules(a.rulesPath)
		if err != nil {
			return Index{}, err
		}
		index = ApplySubscriptionRules(index, a.project, rules)
	}
	index, err = a.governanceIndex(index)
	if err != nil {
		return Index{}, err
	}
	if a.syncIndex != nil {
		if err := a.syncIndex(&index, a.archive); err != nil {
			return Index{}, err
		}
	}
	return index, nil
}

func (a *App) subscriptionRules() (SubscriptionRules, error) {
	if a.loadRules == nil {
		return SubscriptionRules{}, nil
	}
	return a.loadRules(a.rulesPath)
}
