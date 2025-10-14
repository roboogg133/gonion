package torcell

const (
	CellSize = 512
)

const (
	CellTypeVersions      CellType = 7
	CellTypeNetinfo       CellType = 8
	CellTypeCerts         CellType = 129
	CellTypeAuthChallenge CellType = 130
)

type CellType byte

type Cell struct {
	CircuitID uint16
	Type      CellType
	Payload   []byte
}
