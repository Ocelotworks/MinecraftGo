package enums

type VillagerType int

const (
	DesertVillager  VillagerType = 0
	JungleVillager  VillagerType = 1
	PlainsVillager  VillagerType = 2
	SavannaVillager VillagerType = 3
	SnowVillager    VillagerType = 4
	SwampVillager   VillagerType = 5
	TaigaVillager   VillagerType = 6
)

type Pose int

const (
	StandingPose   Pose = 0
	FallFlyingPose Pose = 1
	SleepingPose   Pose = 2
	SwimmingPose   Pose = 3
	SpinAttackPose Pose = 4
	SneakingPose   Pose = 5
	DyingPose      Pose = 6
)

type Direction int

const (
	Down  Direction = 0
	Up    Direction = 1
	North Direction = 2
	South Direction = 3
	West  Direction = 4
	East  Direction = 5
)

type EntityState byte

const (
	OnFire    EntityState = 0x01
	Crouched  EntityState = 0x02
	Unused    EntityState = 0x03
	Sprinting EntityState = 0x04
	Swimming  EntityState = 0x10
	Invisible EntityState = 0x20
	Glowing   EntityState = 0x40
	Elytra    EntityState = 0x80
)
