package dm_xml

import "encoding/xml"

type XMLEstimate struct {
	XMLName         xml.Name `xml:"Document"`
	Text            string   `xml:",chardata" json:"-"`
	ProgramVersion  string   `xml:"ProgramVersion,attr"`
	Generator       string   `xml:"Generator,attr"`
	DocumentType    string   `xml:"DocumentType,attr"`
	GsDocSignatures struct {
		Text string `xml:",chardata" json:"-"`
		Item []struct {
			Text    string `xml:",chardata" json:"-"`
			Caption string `xml:"Caption,attr"`
			ID      string `xml:"ID,attr"`
		} `xml:"Item"`
	} `xml:"GsDocSignatures"`
	RegionalK struct {
		Text    string `xml:",chardata" json:"-"`
		Options string `xml:"Options,attr"`
	} `xml:"RegionalK"`
	TerZoneK struct {
		Text    string `xml:",chardata" json:"-"`
		Caption string `xml:"Caption,attr"`
		Options string `xml:"Options,attr"`
	} `xml:"TerZoneK"`
	Koefficients struct {
		Text string `xml:",chardata" json:"-"`
		K    struct {
			Text    string `xml:",chardata" json:"-"`
			Caption string `xml:"Caption,attr"`
			Code    string `xml:"Code,attr"`
			Options string `xml:"Options,attr"`
			ValueOZ string `xml:"Value_OZ,attr"`
			ValueEM string `xml:"Value_EM,attr"`
			Level   string `xml:"Level,attr"`
		} `xml:"K"`
	} `xml:"Koefficients"`
	WinterCatalog struct {
		Text           string `xml:",chardata" json:"-"`
		Options        string `xml:"Options,attr"`
		WinterMode     string `xml:"WinterMode,attr"`
		WinterLinkMode string `xml:"WinterLinkMode,attr"`
		WinterK        []struct {
			Text    string `xml:",chardata" json:"-"`
			Caption string `xml:"Caption,attr"`
			Code    string `xml:"Code,attr"`
			OZ      string `xml:"OZ,attr"`
			EM      string `xml:"EM,attr"`
			ZM      string `xml:"ZM,attr"`
			MT      string `xml:"MT,attr"`
		} `xml:"WinterK"`
		CommonWinterK string `xml:"CommonWinterK"`
	} `xml:"WinterCatalog"`
	RegionInfo struct {
		Text          string `xml:",chardata" json:"-"`
		RegionName    string `xml:"RegionName,attr"`
		RegionID      string `xml:"RegionID,attr"`
		Zone84Name    string `xml:"Zone84Name,attr"`
		Zone84ID      string `xml:"Zone84ID,attr"`
		Zone01Name    string `xml:"Zone01Name,attr"`
		Zone01ID      string `xml:"Zone01ID,attr"`
		AdmRegionCode string `xml:"AdmRegionCode,attr"`
	} `xml:"RegionInfo"`
	FRSNInfo struct {
		Text      string `xml:",chardata" json:"-"`
		BaseType  string `xml:"BaseType,attr"`
		BaseName  string `xml:"BaseName,attr"`
		RegNumber string `xml:"RegNumber,attr"`
		RegDate   string `xml:"RegDate,attr"`
	} `xml:"FRSN_Info"`
	Parameters struct {
		Text             string `xml:",chardata" json:"-"`
		Options          string `xml:"Options,attr"`
		BasePrices       string `xml:"BasePrices,attr"`
		BaseCalcVrs      string `xml:"BaseCalcVrs,attr"`
		TzDigits         string `xml:"TzDigits,attr"`
		BlockRoundMode   string `xml:"BlockRoundMode,attr"`
		MultKPosCalcMode string `xml:"MultKPosCalcMode,attr"`
		TempZone         string `xml:"TempZone,attr"`
		TsnTempZone      string `xml:"TsnTempZone,attr"`
		MatDigits        string `xml:"MatDigits,attr"`
		MatRoundMode     string `xml:"MatRoundMode,attr"`
		PosKDigits       string `xml:"PosKDigits,attr"`
		ItogOptions      string `xml:"ItogOptions,attr"`
		FirstItogItem    string `xml:"FirstItogItem,attr"`
		ItogExpandTo     string `xml:"ItogExpandTo,attr"`
		PropsConfigName  string `xml:"PropsConfigName,attr"`
		CommonNK         struct {
			Text        string `xml:",chardata" json:"-"`
			ActiveItems string `xml:"ActiveItems,attr"`
		} `xml:"CommonNK"`
		CommonPK struct {
			Text        string `xml:",chardata" json:"-"`
			ActiveItems string `xml:"ActiveItems,attr"`
		} `xml:"CommonPK"`
		MtsnNPZpm struct {
			Text string `xml:",chardata" json:"-"`
			NB   string `xml:"NB,attr"`
			PB   string `xml:"PB,attr"`
			NC   string `xml:"NC,attr"`
			PC   string `xml:"PC,attr"`
		} `xml:"MtsnNPZpm"`
		Numbering struct {
			Text    string `xml:",chardata" json:"-"`
			Mode    string `xml:"Mode,attr"`
			Options string `xml:"Options,attr"`
		} `xml:"Numbering"`
	} `xml:"Parameters"`
	Indexes struct {
		Text            string `xml:",chardata" json:"-"`
		IndexesMode     string `xml:"IndexesMode,attr"`
		IndexesLinkMode string `xml:"IndexesLinkMode,attr"`
		CategoryIndexes struct {
			Text                 string `xml:",chardata" json:"-"`
			Construction         string `xml:"Construction,attr"`
			ConstrTransportation string `xml:"ConstrTransportation,attr"`
		} `xml:"CategoryIndexes"`
		IndexesPos struct {
			Text  string `xml:",chardata" json:"-"`
			Index []struct {
				Text          string `xml:",chardata" json:"-"`
				Caption       string `xml:"Caption,attr"`
				Code          string `xml:"Code,attr"`
				SMR           string `xml:"SMR,attr"`
				OZ            string `xml:"OZ,attr"`
				EM            string `xml:"EM,attr"`
				ZM            string `xml:"ZM,attr"`
				MT            string `xml:"MT,attr"`
				IndexesAddOns struct {
					Text  string `xml:",chardata" json:"-"`
					AddOn []struct {
						Text string `xml:",chardata" json:"-"`
						OZ   string `xml:"OZ,attr"`
						ZM   string `xml:"ZM,attr"`
						Type string `xml:"Type,attr"`
						EM   string `xml:"EM,attr"`
						MT   string `xml:"MT,attr"`
					} `xml:"AddOn"`
				} `xml:"IndexesAddOns"`
			} `xml:"Index"`
		} `xml:"IndexesPos"`
	} `xml:"Indexes"`
	AddZatrats struct {
		Text         string `xml:",chardata" json:"-"`
		AddZatrGlava []struct {
			Text    string `xml:",chardata" json:"-"`
			Glava   string `xml:"Glava,attr"`
			AddZatr struct {
				Text    string `xml:",chardata" json:"-"`
				Caption string `xml:"Caption,attr"`
				Options string `xml:"Options,attr"`
				Formula string `xml:"Formula,attr"`
				Value   string `xml:"Value,attr"`
				Level   string `xml:"Level,attr"`
			} `xml:"AddZatr"`
		} `xml:"AddZatrGlava"`
	} `xml:"AddZatrats"`
	OsInfo struct {
		Text      string `xml:",chardata" json:"-"`
		OSChapter string `xml:"OSChapter,attr"`
		LinkType  string `xml:"LinkType,attr"`
		CCChapter struct {
			Text string `xml:",chardata" json:"-"`
			Cons string `xml:"Cons,attr"`
			Rec  string `xml:"Rec,attr"`
			Road string `xml:"Road,attr"`
		} `xml:"CCChapter"`
	} `xml:"OsInfo"`
	CennikAutoLoad struct {
		Text        string `xml:",chardata" json:"-"`
		Enabled     string `xml:"Enabled,attr"`
		MatchFields string `xml:"MatchFields,attr"`
		Groups      string `xml:"Groups,attr"`
		SubGroups   string `xml:"SubGroups,attr"`
		Options     string `xml:"Options,attr"`
		PriceTypes  string `xml:"PriceTypes,attr"`
		DocLink     string `xml:"DocLink"`
	} `xml:"CennikAutoLoad"`
	VidRabCatalog struct {
		Text    string `xml:",chardata" json:"-"`
		Caption string `xml:"Caption,attr"`
		VidsRab []struct {
			Text        string `xml:",chardata" json:"-"`
			Type        string `xml:"Type,attr"`
			CatFile     string `xml:"CatFile,attr"`
			NrspFile    string `xml:"NrspFile,attr"`
			VidRabGroup []struct {
				Text    string `xml:",chardata" json:"-"`
				Caption string `xml:"Caption,attr"`
				ID      string `xml:"ID,attr"`
				VidRab  []struct {
					Text     string `xml:",chardata" json:"-"`
					Caption  string `xml:"Caption,attr"`
					ID       string `xml:"ID,attr"`
					Nacl     string `xml:"Nacl,attr"`
					Plan     string `xml:"Plan,attr"`
					NaclMask string `xml:"NaclMask,attr"`
					PlanMask string `xml:"PlanMask,attr"`
					OsColumn string `xml:"OsColumn,attr"`
				} `xml:"Vid_Rab"`
			} `xml:"VidRab_Group"`
		} `xml:"Vids_Rab"`
	} `xml:"VidRab_Catalog"`
	Chapters struct {
		Text    string `xml:",chardata" json:"-"`
		Chapter []struct {
			Text     string `xml:",chardata" json:"-"`
			Caption  string `xml:"Caption,attr"`
			SysID    string `xml:"SysID,attr"`
			Position []struct {
				Text         string `xml:",chardata" json:"-"`
				Caption      string `xml:"Caption,attr"`
				Number       string `xml:"Number,attr"`
				Code         string `xml:"Code,attr"`
				ColorIndex   string `xml:"ColorIndex,attr"`
				Units        string `xml:"Units,attr"`
				SysID        string `xml:"SysID,attr"`
				PriceLevel   string `xml:"PriceLevel,attr"`
				DBComment    string `xml:"DBComment,attr"`
				AttrQuantity string `xml:"Quantity,attr"`
				IndexCode    string `xml:"IndexCode,attr"`
				DBFlags      string `xml:"DBFlags,attr"`
				PzSync       string `xml:"PzSync,attr"`
				Vr2001       string `xml:"Vr2001,attr"`
				WinterK      string `xml:"WinterK,attr"`
				Mass         string `xml:"Mass,attr"`
				Quantity     struct {
					Text      string `xml:",chardata" json:"-"`
					Fx        string `xml:"Fx,attr"`
					KUnit     string `xml:"KUnit,attr"`
					Precision string `xml:"Precision,attr"`
					Result    string `xml:"Result,attr"`
				} `xml:"Quantity"`
				PriceBase struct {
					Text string `xml:",chardata" json:"-"`
					PZ   string `xml:"PZ,attr"`
					OZ   string `xml:"OZ,attr"`
					EM   string `xml:"EM,attr"`
					ZM   string `xml:"ZM,attr"`
					MT   string `xml:"MT,attr"`
				} `xml:"PriceBase"`
				Resources struct {
					Text string `xml:",chardata" json:"-"`
					Tzr  struct {
						Text     string `xml:",chardata" json:"-"`
						Caption  string `xml:"Caption,attr"`
						Code     string `xml:"Code,attr"`
						Units    string `xml:"Units,attr"`
						Quantity string `xml:"Quantity,attr"`
					} `xml:"Tzr"`
					Mch []struct {
						Text      string `xml:",chardata" json:"-"`
						Caption   string `xml:"Caption,attr"`
						Code      string `xml:"Code,attr"`
						Units     string `xml:"Units,attr"`
						Quantity  string `xml:"Quantity,attr"`
						PriceBase struct {
							Text  string `xml:",chardata" json:"-"`
							Value string `xml:"Value,attr"`
							ZM    string `xml:"ZM,attr"`
						} `xml:"PriceBase"`
					} `xml:"Mch"`
					Mat []struct {
						Text      string `xml:",chardata" json:"-"`
						Caption   string `xml:"Caption,attr"`
						Code      string `xml:"Code,attr"`
						Units     string `xml:"Units,attr"`
						Quantity  string `xml:"Quantity,attr"`
						Mass      string `xml:"Mass,attr"`
						Options   string `xml:"Options,attr"`
						Attribs   string `xml:"Attribs,attr"`
						PriceBase struct {
							Text  string `xml:",chardata" json:"-"`
							Value string `xml:"Value,attr"`
						} `xml:"PriceBase"`
					} `xml:"Mat"`
				} `xml:"Resources"`
				WorksList struct {
					Text string `xml:",chardata" json:"-"`
					Work []struct {
						Text    string `xml:",chardata" json:"-"`
						Caption string `xml:"Caption,attr"`
					} `xml:"Work"`
				} `xml:"WorksList"`
				Koefficients struct {
					Text string `xml:",chardata" json:"-"`
					K    struct {
						Text    string `xml:",chardata" json:"-"`
						Caption string `xml:"Caption,attr"`
						Options string `xml:"Options,attr"`
						ValuePZ string `xml:"Value_PZ,attr"`
						Level   string `xml:"Level,attr"`
					} `xml:"K"`
				} `xml:"Koefficients"`
			} `xml:"Position"`
		} `xml:"Chapter"`
	} `xml:"Chapters"`
	ReportOptions struct {
		Text          string `xml:",chardata" json:"-"`
		Options       string `xml:"Options,attr"`
		Kk            string `xml:"Kk,attr"`
		RangingGroups string `xml:"RangingGroups,attr"`
		RangingRates  struct {
			Text string `xml:",chardata" json:"-"`
			Mch  string `xml:"Mch,attr"`
			Mat  string `xml:"Mat,attr"`
		} `xml:"RangingRates"`
	} `xml:"ReportOptions"`
}
