package datasets

type SocketType Enum

type ConduitRank struct {
	Identifiables
	Tier           int           `json:"tier"`
	SpellTooltipID int           ``
	SpellTooltip   *SpellTooltip `json:"spell_tooltip" pg:"rel:has-one"`
}

type Conduit struct {
	Identifiable
	Name       LocalizedField `json:"name"`
	ItemID     int            ``
	Item       *Item          `json:"item" pg:"rel:has-one"`
	SocketType *SocketType    `json:"socket_type"`
	Ranks      []*ConduitRank `json:"ranks"`
}

type Covenant struct {
	ID                 int                `json:"id"`
	Name               LocalizedField     `json:"name"`
	Description        LocalizedField     `json:"description"`
	SignatureAbilityID int                ``
	SignatureAbility   *CovenantAbility   `json:"signature_ability" pg:"-"`
	ClassAbilities     []*CovenantAbility `json:"class_abilities" pg:"-"`
	Soulbinds          []*Soulbind        `json:"soulbinds" pg:"-"`
	RenownRewards      []*RenownRewards   `json:"renown_rewards" pg:"-"`
}

type CovenantAbility struct {
	Identifiable
	CovenantID      int            ``
	Covenant        *Covenant      `pg:"rel:has-one"`
	PlayableClassID int            ``
	PlayableClass   *PlayableClass `json:"playable_class" pg:"rel:has-one"`
	SpellTooltipID  int            ``
	SpellTooltip    *SpellTooltip  `json:"spell_tooltip" pg:"rel:has-one"`
}

type Reward NamedItem

type RenownRewards []struct {
	CovenantID int       ``
	Covenant   *Covenant `pg:"rel:has-one"`
	RewardID   int       ``
	Reward     *Reward   `json:"reward" pg:"rel:has-one"`
	Level      int       `json:"level"`
}

type Follower NamedItem

type Soulbind struct {
	Identifiable
	Name             LocalizedField  `json:"name"`
	CovenantID       int             ``
	Covenant         *Covenant       `json:"covenant" pg:"rel:has-one"`
	CreatureID       int             ``
	Creature         *Creature       `json:"creature" pg:"rel:has-one"`
	FollowerID       int             ``
	Follower         *Follower       `json:"follower" pg:"rel:has-one"`
	TechTalentTreeID int             ``
	TechTalentTree   *TechTalentTree `json:"talent_tree" pg:"rel:has-one"`
}