package staruml

type Ref struct {
	Ref string `json:"$ref"`
}

type Element struct {
	ElemType      string      `json:"_type"`
	ID            string      `json:"_id"`
	Parent        *Ref        `json:"_parent"`
	Name          string      `json:"name"`
	Length        interface{} `json:"length"`
	Type          string      `json:"type"`
	PrimaryKey    bool        `json:"primaryKey"`
	Unique        bool        `json:"unique"`
	Nullable      bool        `json:"nullable"`
	Documentation string      `json:"documentation"`
	Head          *Ref        `json:"head"`
	End1          *Ref        `json:"end1"`
	End2          *Ref        `json:"end2"`
	Reference     *Ref        `json:"reference"`
	ReferenceTo   *Ref        `json:"referenceTo"`

	OwnedElements []*Element `json:"ownedElements"`
	Columns       []*Element `json:"columns"`
}
