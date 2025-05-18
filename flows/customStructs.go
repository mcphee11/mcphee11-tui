package flows

import "time"

// ------------------------------ CUSTOM STRUCT UNTIL SDK UPDATED ---------------------------

type FlowentitylistingCUSTOM struct {
	// Entities
	Entities *[]FlowCUSTOM `json:"entities,omitempty"`
	// PageSize
	PageSize *int `json:"pageSize,omitempty"`
	// PageNumber
	PageNumber *int `json:"pageNumber,omitempty"`
	// Total
	Total *int `json:"total,omitempty"`
	// FirstUri
	FirstUri *string `json:"firstUri,omitempty"`
	// SelfUri
	SelfUri *string `json:"selfUri,omitempty"`
	// NextUri
	NextUri *string `json:"nextUri,omitempty"`
	// PreviousUri
	PreviousUri *string `json:"previousUri,omitempty"`
	// LastUri
	LastUri *string `json:"lastUri,omitempty"`
	// PageCount
	PageCount *int `json:"pageCount,omitempty"`
}

type FlowCUSTOM struct {
	// Id - The flow identifier
	Id *string `json:"id,omitempty"`
	// Name - The flow name
	Name *string `json:"name,omitempty"`
	// Division - The division to which this entity belongs.
	Division *WritabledivisionCUSTOM `json:"division,omitempty"`
	// Description
	Description *string `json:"description,omitempty"`
	// VarType
	VarType *string `json:"type,omitempty"`
	// LockedUser - User that has the flow locked.
	LockedUser *UserCUSTOM `json:"lockedUser,omitempty"`
	// LockedClient - OAuth client that has the flow locked.
	LockedClient *DomainentityrefCUSTOM `json:"lockedClient,omitempty"`
	// Active
	Active *bool `json:"active,omitempty"`
	// System
	System *bool `json:"system,omitempty"`
	// Deleted
	Deleted *bool `json:"deleted,omitempty"`
	// PublishedVersion
	PublishedVersion *FlowversionCUSTOM `json:"publishedVersion,omitempty"`
	// SavedVersion
	SavedVersion *FlowversionCUSTOM `json:"savedVersion,omitempty"`
	// InputSchema - json schema describing the inputs for the flow
	InputSchema *map[string]interface{} `json:"inputSchema,omitempty"`
	// OutputSchema - json schema describing the outputs for the flow
	OutputSchema *map[string]interface{} `json:"outputSchema,omitempty"`
	// CheckedInVersion
	CheckedInVersion *FlowversionCUSTOM `json:"checkedInVersion,omitempty"`
	// DebugVersion
	DebugVersion *FlowversionCUSTOM `json:"debugVersion,omitempty"`
	// PublishedBy
	PublishedBy *UserCUSTOM `json:"publishedBy,omitempty"`
	// CurrentOperation
	CurrentOperation *OperationCUSTOM `json:"currentOperation,omitempty"`
	// NluInfo - Information about the natural language understanding configuration for the published version of the flow
	NluInfo *NluinfoCUSTOM `json:"nluInfo,omitempty"`
	// SupportedLanguages - List of supported languages for the published version of the flow.
	SupportedLanguages *[]SupportedlanguageCUSTOM `json:"supportedLanguages,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type WritabledivisionCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type DomainentityrefCUSTOM struct {
	// Id
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// SelfUri
	SelfUri *string `json:"selfUri,omitempty"`
}

type UserCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Division - The division to which this entity belongs.
	Division *DivisionCUSTOM `json:"division,omitempty"`
	// Chat
	Chat *ChatCUSTOM `json:"chat,omitempty"`
	// Department
	Department *string `json:"department,omitempty"`
	// Email
	Email *string `json:"email,omitempty"`
	// PrimaryContactInfo - Auto populated from addresses.
	PrimaryContactInfo *[]ContactCUSTOM `json:"primaryContactInfo,omitempty"`
	// Addresses - Email addresses and phone numbers for this user
	Addresses *[]ContactCUSTOM `json:"addresses,omitempty"`
	// State - The current state for this user.
	State *string `json:"state,omitempty"`
	// Title
	Title *string `json:"title,omitempty"`
	// Username
	Username *string `json:"username,omitempty"`
	// Manager
	Manager **UserCUSTOM `json:"manager,omitempty"`
	// Images
	Images *[]UserimageCUSTOM `json:"images,omitempty"`
	// Version - Required when updating a user, this value should be the current version of the user.  The current version can be obtained with a GET on the user before doing a PATCH.
	Version *int `json:"version,omitempty"`
	// Certifications
	Certifications *[]string `json:"certifications,omitempty"`
	// Biography
	Biography *BiographyCUSTOM `json:"biography,omitempty"`
	// EmployerInfo
	EmployerInfo *EmployerinfoCUSTOM `json:"employerInfo,omitempty"`
	// RoutingStatus - ACD routing status
	RoutingStatus *RoutingstatusCUSTOM `json:"routingStatus,omitempty"`
	// Presence - Active presence
	Presence *UserpresenceCUSTOM `json:"presence,omitempty"`
	// ConversationSummary - Summary of conversion statistics for conversation types.
	ConversationSummary *UserconversationsummaryCUSTOM `json:"conversationSummary,omitempty"`
	// OutOfOffice - Determine if out of office is enabled
	OutOfOffice **OutofofficeCUSTOM `json:"outOfOffice,omitempty"`
	// Geolocation - Current geolocation position
	Geolocation *GeolocationCUSTOM `json:"geolocation,omitempty"`
	// Station - Effective, default, and last station information
	Station **UserstationsCUSTOM `json:"station,omitempty"`
	// Authorization - Roles and permissions assigned to the user
	Authorization *UserauthorizationCUSTOM `json:"authorization,omitempty"`
	// ProfileSkills - Profile skills possessed by the user
	ProfileSkills *[]string `json:"profileSkills,omitempty"`
	// Locations - The user placement at each site location.
	Locations *[]LocationCUSTOM `json:"locations,omitempty"`
	// Groups - The groups the user is a member of
	Groups *[]GroupCUSTOM `json:"groups,omitempty"`
	// Team - The team the user is a member of
	Team *TeamCUSTOM `json:"team,omitempty"`
	// Skills - Routing (ACD) skills possessed by the user
	Skills *[]UserroutingskillCUSTOM `json:"skills,omitempty"`
	// Languages - Routing (ACD) languages possessed by the user
	Languages *[]UserroutinglanguageCUSTOM `json:"languages,omitempty"`
	// AcdAutoAnswer - acd auto answer
	AcdAutoAnswer *bool `json:"acdAutoAnswer,omitempty"`
	// LanguagePreference - preferred language by the user
	LanguagePreference *string `json:"languagePreference,omitempty"`
	// LastTokenIssued
	LastTokenIssued *OauthlasttokenissuedCUSTOM `json:"lastTokenIssued,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type DivisionCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type ChatCUSTOM struct {
	// JabberId
	JabberId *string `json:"jabberId,omitempty"`
}

type ContactCUSTOM struct {
	// Address - Email address or phone number for this contact type
	Address *string `json:"address,omitempty"`
	// Display - Formatted version of the address property
	Display *string `json:"display,omitempty"`
	// MediaType
	MediaType *string `json:"mediaType,omitempty"`
	// VarType
	VarType *string `json:"type,omitempty"`
	// Extension - Use internal extension instead of address. Mutually exclusive with the address field.
	Extension *string `json:"extension,omitempty"`
	// CountryCode
	CountryCode *string `json:"countryCode,omitempty"`
}

type UserimageCUSTOM struct {
	// Resolution - Height and/or width of image. ex: 640x480 or x128
	Resolution *string `json:"resolution,omitempty"`
	// ImageUri
	ImageUri *string `json:"imageUri,omitempty"`
}

type BiographyCUSTOM struct {
	// Biography - Personal detailed description
	Biography *string `json:"biography,omitempty"`
	// Interests
	Interests *[]string `json:"interests,omitempty"`
	// Hobbies
	Hobbies *[]string `json:"hobbies,omitempty"`
	// Spouse
	Spouse *string `json:"spouse,omitempty"`
	// Education - User education details
	Education *[]EducationCUSTOM `json:"education,omitempty"`
}

type EducationCUSTOM struct {
	// School
	School *string `json:"school,omitempty"`
	// FieldOfStudy
	FieldOfStudy *string `json:"fieldOfStudy,omitempty"`
	// Notes - Notes about education has a 2000 character limit
	Notes *string `json:"notes,omitempty"`
	// DateStart - Dates are represented as an ISO-8601 string. For example: yyyy-MM-dd
	DateStart *time.Time `json:"dateStart,omitempty"`
	// DateEnd - Dates are represented as an ISO-8601 string. For example: yyyy-MM-dd
	DateEnd *time.Time `json:"dateEnd,omitempty"`
}

type EmployerinfoCUSTOM struct {
	// OfficialName
	OfficialName *string `json:"officialName,omitempty"`
	// EmployeeId
	EmployeeId *string `json:"employeeId,omitempty"`
	// EmployeeType
	EmployeeType *string `json:"employeeType,omitempty"`
	// DateHire
	DateHire *string `json:"dateHire,omitempty"`
}

type RoutingstatusCUSTOM struct {
	// UserId - The userId of the agent
	UserId *string `json:"userId,omitempty"`
	// Status - Indicates the Routing State of the agent.  A value of OFF_QUEUE will be returned if the specified user does not exist.
	Status *string `json:"status,omitempty"`
	// StartTime - The timestamp when the agent went into this state. Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	StartTime *time.Time `json:"startTime,omitempty"`
}

type UserpresenceCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Source - Represents the source where the Presence was set. Some examples are: PURECLOUD, LYNC, OUTLOOK, etc.
	Source *string `json:"source,omitempty"`
	// Primary - A boolean used to tell whether or not to set this presence source as the primary on a PATCH
	Primary *bool `json:"primary,omitempty"`
	// PresenceDefinition
	PresenceDefinition *PresencedefinitionCUSTOM `json:"presenceDefinition,omitempty"`
	// Message
	Message *string `json:"message,omitempty"`
	// ModifiedDate - Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	ModifiedDate *time.Time `json:"modifiedDate,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type PresencedefinitionCUSTOM struct {
	// Id - description
	Id *string `json:"id,omitempty"`
	// SystemPresence
	SystemPresence *string `json:"systemPresence,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type UserconversationsummaryCUSTOM struct {
	// UserId
	UserId *string `json:"userId,omitempty"`
	// Call
	Call *MediasummaryCUSTOM `json:"call,omitempty"`
	// Callback
	Callback *MediasummaryCUSTOM `json:"callback,omitempty"`
	// Email
	Email *MediasummaryCUSTOM `json:"email,omitempty"`
	// Message
	Message *MediasummaryCUSTOM `json:"message,omitempty"`
	// Chat
	Chat *MediasummaryCUSTOM `json:"chat,omitempty"`
	// SocialExpression
	SocialExpression *MediasummaryCUSTOM `json:"socialExpression,omitempty"`
	// Video
	Video *MediasummaryCUSTOM `json:"video,omitempty"`
}

type MediasummaryCUSTOM struct {
	// ContactCenter
	ContactCenter *MediasummarydetailCUSTOM `json:"contactCenter,omitempty"`

	// Enterprise
	Enterprise *MediasummarydetailCUSTOM `json:"enterprise,omitempty"`
}

type MediasummarydetailCUSTOM struct {
	// Active
	Active *int `json:"active,omitempty"`
	// Acw
	Acw *int `json:"acw,omitempty"`
}

type OutofofficeCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// User
	User **UserCUSTOM `json:"user,omitempty"`
	// StartDate - Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	StartDate *time.Time `json:"startDate,omitempty"`
	// EndDate - Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	EndDate *time.Time `json:"endDate,omitempty"`
	// Active
	Active *bool `json:"active,omitempty"`
	// Indefinite
	Indefinite *bool `json:"indefinite,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type GeolocationCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// VarType - A string used to describe the type of client the geolocation is being updated from e.g. ios, android, web, etc.
	VarType *string `json:"type,omitempty"`
	// Primary - A boolean used to tell whether or not to set this geolocation client as the primary on a PATCH
	Primary *bool `json:"primary,omitempty"`
	// Latitude
	Latitude *float64 `json:"latitude,omitempty"`
	// Longitude
	Longitude *float64 `json:"longitude,omitempty"`
	// Country
	Country *string `json:"country,omitempty"`
	// Region
	Region *string `json:"region,omitempty"`
	// City
	City *string `json:"city,omitempty"`
	// Locations
	Locations *[]LocationdefinitionCUSTOM `json:"locations,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type LocationdefinitionCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// ContactUser - Site contact for the location entity
	ContactUser *AddressableentityrefCUSTOM `json:"contactUser,omitempty"`
	// EmergencyNumber - Emergency number for the location entity
	EmergencyNumber *LocationemergencynumberCUSTOM `json:"emergencyNumber,omitempty"`
	// Address
	Address *LocationaddressCUSTOM `json:"address,omitempty"`
	// State - Current state of the location entity
	State *string `json:"state,omitempty"`
	// Notes - Notes for the location entity
	Notes *string `json:"notes,omitempty"`
	// Version - Current version of the location entity, value to be supplied should be retrieved by a GET or on create/update response
	Version *int `json:"version,omitempty"`
	// Path - A list of ancestor IDs in order
	Path *[]string `json:"path,omitempty"`
	// ProfileImage - Profile image of the location entity, retrieved with ?expand=images query parameter
	ProfileImage *[]LocationimageCUSTOM `json:"profileImage,omitempty"`
	// FloorplanImage - Floorplan images of the location entity, retrieved with ?expand=images query parameter
	FloorplanImage *[]LocationimageCUSTOM `json:"floorplanImage,omitempty"`
	// AddressVerificationDetails - Address verification information, retrieve dwith the ?expand=addressVerificationDetails query parameter
	AddressVerificationDetails *LocationaddressverificationdetailsCUSTOM `json:"addressVerificationDetails,omitempty"`
	// AddressVerified - Boolean field which states if the address has been verified as an actual address
	AddressVerified *bool `json:"addressVerified,omitempty"`
	// AddressStored - Boolean field which states if the address has been stored for E911
	AddressStored *bool `json:"addressStored,omitempty"`
	// Images
	Images *string `json:"images,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type AddressableentityrefCUSTOM struct {
	// Id
	Id *string `json:"id,omitempty"`
	// SelfUri
	SelfUri *string `json:"selfUri,omitempty"`
}

type LocationemergencynumberCUSTOM struct {
	// E164
	E164 *string `json:"e164,omitempty"`
	// Number
	Number *string `json:"number,omitempty"`
	// VarType - The type of emergency number.
	VarType *string `json:"type,omitempty"`
}

type LocationaddressCUSTOM struct {
	// City
	City *string `json:"city,omitempty"`
	// Country
	Country *string `json:"country,omitempty"`
	// CountryName
	CountryName *string `json:"countryName,omitempty"`
	// State
	State *string `json:"state,omitempty"`
	// Street1
	Street1 *string `json:"street1,omitempty"`
	// Street2
	Street2 *string `json:"street2,omitempty"`
	// Zipcode
	Zipcode *string `json:"zipcode,omitempty"`
}

type LocationimageCUSTOM struct {
	// Resolution - Height and/or width of image. ex: 640x480 or x128
	Resolution *string `json:"resolution,omitempty"`

	// ImageUri
	ImageUri *string `json:"imageUri,omitempty"`
}

type LocationaddressverificationdetailsCUSTOM struct {
	// Status - Status of address verification process
	Status *string `json:"status,omitempty"`

	// DateFinished - Finished time of address verification process. Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	DateFinished *time.Time `json:"dateFinished,omitempty"`

	// DateStarted - Time started of address verification process. Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	DateStarted *time.Time `json:"dateStarted,omitempty"`

	// Service - Third party service used for address verification
	Service *string `json:"service,omitempty"`
}

type UserstationsCUSTOM struct {
	// AssociatedStation - Current associated station for this user.
	AssociatedStation *UserstationCUSTOM `json:"associatedStation,omitempty"`
	// EffectiveStation - The station where the user can be reached based on their default and associated station.
	EffectiveStation *UserstationCUSTOM `json:"effectiveStation,omitempty"`
	// DefaultStation - Default station to be used if not associated with a station.
	DefaultStation *UserstationCUSTOM `json:"defaultStation,omitempty"`
	// LastAssociatedStation - Last associated station for this user.
	LastAssociatedStation *UserstationCUSTOM `json:"lastAssociatedStation,omitempty"`
}

type UserstationCUSTOM struct {
	// Id - A globally unique identifier for this station
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// VarType
	VarType *string `json:"type,omitempty"`
	// AssociatedUser
	AssociatedUser *UserCUSTOM `json:"associatedUser,omitempty"`
	// AssociatedDate - Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	AssociatedDate *time.Time `json:"associatedDate,omitempty"`
	// DefaultUser
	DefaultUser *UserCUSTOM `json:"defaultUser,omitempty"`
	// ProviderInfo - Provider-specific info for this station, e.g. { \"edgeGroupId\": \"ffe7b15c-a9cc-4f4c-88f5-781327819a49\" }
	ProviderInfo *map[string]string `json:"providerInfo,omitempty"`
}

type UserauthorizationCUSTOM struct {
	// Roles
	Roles *[]DomainroleCUSTOM `json:"roles,omitempty"`
	// UnusedRoles - A collection of the roles the user is not using
	UnusedRoles *[]DomainroleCUSTOM `json:"unusedRoles,omitempty"`
	// Permissions - A collection of the permissions granted by all assigned roles
	Permissions *[]string `json:"permissions,omitempty"`
	// PermissionPolicies - The policies configured for assigned permissions.
	PermissionPolicies *[]ResourcepermissionpolicyCUSTOM `json:"permissionPolicies,omitempty"`
}

type DomainroleCUSTOM struct {
	// Id - The ID of the role
	Id *string `json:"id,omitempty"`
	// Name - The name of the role
	Name *string `json:"name,omitempty"`
}

type ResourcepermissionpolicyCUSTOM struct {
	// Id
	Id *string `json:"id,omitempty"`
	// Domain
	Domain *string `json:"domain,omitempty"`
	// EntityName
	EntityName *string `json:"entityName,omitempty"`
	// PolicyName
	PolicyName *string `json:"policyName,omitempty"`
	// PolicyDescription
	PolicyDescription *string `json:"policyDescription,omitempty"`
	// ActionSetKey
	ActionSetKey *string `json:"actionSetKey,omitempty"`
	// AllowConditions
	AllowConditions *bool `json:"allowConditions,omitempty"`
	// ResourceConditionNode
	ResourceConditionNode *ResourceconditionnodeCUSTOM `json:"resourceConditionNode,omitempty"`
	// NamedResources
	NamedResources *[]string `json:"namedResources,omitempty"`
	// ResourceCondition
	ResourceCondition *string `json:"resourceCondition,omitempty"`
	// ActionSet
	ActionSet *[]string `json:"actionSet,omitempty"`
}

type ResourceconditionnodeCUSTOM struct {
	// VariableName
	VariableName *string `json:"variableName,omitempty"`
	// Conjunction
	Conjunction *string `json:"conjunction,omitempty"`
	// Operator
	Operator *string `json:"operator,omitempty"`
	// Operands
	Operands *[]ResourceconditionvalueCUSTOM `json:"operands,omitempty"`
	// Terms
	Terms *[]ResourceconditionnodeCUSTOM `json:"terms,omitempty"`
}

type ResourceconditionvalueCUSTOM struct {
	// VarType
	VarType *string `json:"type,omitempty"`
	// Value
	Value *string `json:"value,omitempty"`
}

type LocationCUSTOM struct {
	// Id - Unique identifier for the location
	Id *string `json:"id,omitempty"`
	// FloorplanId - Unique identifier for the location floorplan image
	FloorplanId *string `json:"floorplanId,omitempty"`
	// Coordinates - Users coordinates on the floorplan. Only used when floorplanImage is set
	Coordinates *map[string]float64 `json:"coordinates,omitempty"`
	// Notes - Optional description on the users location
	Notes *string `json:"notes,omitempty"`
	// LocationDefinition
	LocationDefinition *LocationdefinitionCUSTOM `json:"locationDefinition,omitempty"`
}

type GroupCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name - The group name.
	Name *string `json:"name,omitempty"`
	// Description
	Description *string `json:"description,omitempty"`
	// DateModified - Last modified date/time. Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	DateModified *time.Time `json:"dateModified,omitempty"`
	// MemberCount - Number of members.
	MemberCount *int `json:"memberCount,omitempty"`
	// State - Active, inactive, or deleted state.
	State *string `json:"state,omitempty"`
	// Version - Current version for this resource.
	Version *int `json:"version,omitempty"`
	// VarType - Type of group.
	VarType *string `json:"type,omitempty"`
	// Images
	Images *[]UserimageCUSTOM `json:"images,omitempty"`
	// Addresses
	Addresses *[]GroupcontactCUSTOM `json:"addresses,omitempty"`
	// RulesVisible - Are membership rules visible to the person requesting to view the group
	RulesVisible *bool `json:"rulesVisible,omitempty"`
	// Visibility - Who can view this group
	Visibility *string `json:"visibility,omitempty"`
	// Owners - Owners of the group
	Owners *[]UserCUSTOM `json:"owners,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type GroupcontactCUSTOM struct {
	// Address - Phone number for this contact type
	Address *string `json:"address,omitempty"`
	// Extension - Extension is set if the number is e164 valid
	Extension *string `json:"extension,omitempty"`
	// Display - Formatted version of the address property
	Display *string `json:"display,omitempty"`
	// VarType - Contact type of the address
	VarType *string `json:"type,omitempty"`
	// MediaType - Media type of the address
	MediaType *string `json:"mediaType,omitempty"`
}

type TeamCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name - The team name
	Name *string `json:"name,omitempty"`
	// Description - Team information.
	Description *string `json:"description,omitempty"`
	// DateModified - Last modified datetime. Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	DateModified *time.Time `json:"dateModified,omitempty"`
	// MemberCount - Number of members in a team
	MemberCount *int `json:"memberCount,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type UserroutingskillCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Proficiency - A rating from 0.0 to 5.0 that indicates how adept an agent is at a particular skill. When \"Best available skills\" is enabled for a queue in Genesys Cloud, ACD interactions in that queue are routed to agents with higher proficiency ratings.
	Proficiency *float64 `json:"proficiency,omitempty"`
	// State - Activate or deactivate this routing skill.
	State *string `json:"state,omitempty"`
	// SkillUri - URI to the organization skill used by this user skill.
	SkillUri *string `json:"skillUri,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type UserroutinglanguageCUSTOM struct {
	// Id - The globally unique identifier for the object.
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// Proficiency - A rating from 0.0 to 5.0 that indicates how fluent an agent is in a particular language. ACD interactions are routed to agents with higher proficiency ratings.
	Proficiency *float64 `json:"proficiency,omitempty"`
	// State - Activate or deactivate this routing language.
	State *string `json:"state,omitempty"`
	// LanguageUri - URI to the organization language used by this user language.
	LanguageUri *string `json:"languageUri,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type OauthlasttokenissuedCUSTOM struct {
	// DateIssued - Date time is represented as an ISO-8601 string. For example: yyyy-MM-ddTHH:mm:ss[.mmm]Z
	DateIssued *time.Time `json:"dateIssued,omitempty"`
}

type FlowversionCUSTOM struct {
	// Id - The flow version identifier
	Id *string `json:"id,omitempty"`
	// Name
	Name *string `json:"name,omitempty"`
	// CommitVersion
	CommitVersion *string `json:"commitVersion,omitempty"`
	// ConfigurationVersion
	ConfigurationVersion *string `json:"configurationVersion,omitempty"`
	// VarType
	VarType *string `json:"type,omitempty"`
	// Secure
	Secure *bool `json:"secure,omitempty"`
	// Debug
	Debug *bool `json:"debug,omitempty"`
	// CreatedBy
	CreatedBy *UserCUSTOM `json:"createdBy,omitempty"`
	// CreatedByClient
	CreatedByClient *DomainentityrefCUSTOM `json:"createdByClient,omitempty"`
	// ConfigurationUri
	ConfigurationUri *string `json:"configurationUri,omitempty"`
	// DateCreated
	DateCreated *int `json:"dateCreated,omitempty"`
	// GenerationId
	GenerationId *string `json:"generationId,omitempty"`
	// PublishResultUri
	PublishResultUri *string `json:"publishResultUri,omitempty"`
	// InputSchema
	InputSchema *JsonschemadocumentCUSTOM `json:"inputSchema,omitempty"`
	// OutputSchema
	OutputSchema *JsonschemadocumentCUSTOM `json:"outputSchema,omitempty"`
	// NluInfo - Information about the natural language understanding configuration for the flow version
	NluInfo *NluinfoCUSTOM `json:"nluInfo,omitempty"`
	// SupportedLanguages - List of supported languages for this version of the flow
	SupportedLanguages *[]SupportedlanguageCUSTOM `json:"supportedLanguages,omitempty"`
	// SelfUri - The URI for this object
	SelfUri *string `json:"selfUri,omitempty"`
}

type JsonschemadocumentCUSTOM struct {
	// Id
	Id *string `json:"id,omitempty"`
	// Schema
	Schema *string `json:"$schema,omitempty"`
	// Title
	Title *string `json:"title,omitempty"`
	// Description
	Description *string `json:"description,omitempty"`
	// VarType
	VarType *string `json:"type,omitempty"`
	// Required
	Required *[]string `json:"required,omitempty"`
	// Properties
	Properties *map[string]interface{} `json:"properties,omitempty"`
	// AdditionalProperties
	//AdditionalProperties *map[string]interface{} `json:"additionalProperties,omitempty"`
	AdditionalProperties *bool `json:"additionalProperties,omitempty"`
}

type NluinfoCUSTOM struct {
	// Intents
	Intents *[]IntentCUSTOM `json:"intents,omitempty"`
}

type IntentCUSTOM struct {
	// Name
	Name *string `json:"name,omitempty"`
}

type SupportedlanguageCUSTOM struct {
	// Language - Architect supported language tag, e.g. en-us, es-us
	Language *string `json:"language,omitempty"`
	// IsDefault - Whether or not this language is the default language
	IsDefault *bool `json:"isDefault,omitempty"`
}

type OperationCUSTOM struct {
	// Id
	Id *string `json:"id,omitempty"`
	// Complete
	Complete *bool `json:"complete,omitempty"`
	// User
	User *UserCUSTOM `json:"user,omitempty"`
	// Client
	Client *DomainentityrefCUSTOM `json:"client,omitempty"`
	// ErrorMessage
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// ErrorCode
	ErrorCode *string `json:"errorCode,omitempty"`
	// ErrorDetails
	ErrorDetails *[]DetailCUSTOM `json:"errorDetails,omitempty"`
	// ErrorMessageParams
	ErrorMessageParams *map[string]string `json:"errorMessageParams,omitempty"`
	// ActionName - Action name
	ActionName *string `json:"actionName,omitempty"`
	// ActionStatus - Action status
	ActionStatus *string `json:"actionStatus,omitempty"`
}

type DetailCUSTOM struct {
	// ErrorCode
	ErrorCode *string `json:"errorCode,omitempty"`
	// FieldName
	FieldName *string `json:"fieldName,omitempty"`
	// EntityId
	EntityId *string `json:"entityId,omitempty"`
	// EntityName
	EntityName *string `json:"entityName,omitempty"`
}
