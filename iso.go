package kensho

import (
	"golang.org/x/text/language"
)

type (
	iso3166 struct {
		CountryName string
		Alpha2      string
		Alpha3      string
		Numeric     string
	}
)

var iso3166List = []iso3166{
	{"Afghanistan", "AF", "AFG", "004"},
	{"Albania", "AL", "ALB", "008"},
	{"Antarctica", "AQ", "ATA", "010"},
	{"Algeria", "DZ", "DZA", "012"},
	{"American Samoa", "AS", "ASM", "016"},
	{"Andorra", "AD", "AND", "020"},
	{"Angola", "AO", "AGO", "024"},
	{"Antigua and Barbuda", "AG", "ATG", "028"},
	{"Azerbaijan", "AZ", "AZE", "031"},
	{"Argentina", "AR", "ARG", "032"},
	{"Australia", "AU", "AUS", "036"},
	{"Austria", "AT", "AUT", "040"},
	{"Bahamas (the)", "BS", "BHS", "044"},
	{"Bahrain", "BH", "BHR", "048"},
	{"Bangladesh", "BD", "BGD", "050"},
	{"Armenia", "AM", "ARM", "051"},
	{"Barbados", "BB", "BRB", "052"},
	{"Belgium", "BE", "BEL", "056"},
	{"Bermuda", "BM", "BMU", "060"},
	{"Bhutan", "BT", "BTN", "064"},
	{"Bolivia (Plurinational State of)", "BO", "BOL", "068"},
	{"Bosnia and Herzegovina", "BA", "BIH", "070"},
	{"Botswana", "BW", "BWA", "072"},
	{"Bouvet Island", "BV", "BVT", "074"},
	{"Brazil", "BR", "BRA", "076"},
	{"Belize", "BZ", "BLZ", "084"},
	{"British Indian Ocean Territory (the)", "IO", "IOT", "086"},
	{"Solomon Islands", "SB", "SLB", "090"},
	{"Virgin Islands (British)", "VG", "VGB", "092"},
	{"Brunei Darussalam", "BN", "BRN", "096"},
	{"Bulgaria", "BG", "BGR", "100"},
	{"Myanmar", "MM", "MMR", "104"},
	{"Burundi", "BI", "BDI", "108"},
	{"Belarus", "BY", "BLR", "112"},
	{"Cambodia", "KH", "KHM", "116"},
	{"Cameroon", "CM", "CMR", "120"},
	{"Canada", "CA", "CAN", "124"},
	{"Cabo Verde", "CV", "CPV", "132"},
	{"Cayman Islands (the)", "KY", "CYM", "136"},
	{"Central African Republic (the)", "CF", "CAF", "140"},
	{"Sri Lanka", "LK", "LKA", "144"},
	{"Chad", "TD", "TCD", "148"},
	{"Chile", "CL", "CHL", "152"},
	{"China", "CN", "CHN", "156"},
	{"Taiwan (Province of China)", "TW", "TWN", "158"},
	{"Christmas Island", "CX", "CXR", "162"},
	{"Cocos (Keeling) Islands (the)", "CC", "CCK", "166"},
	{"Colombia", "CO", "COL", "170"},
	{"Comoros (the)", "KM", "COM", "174"},
	{"Mayotte", "YT", "MYT", "175"},
	{"Congo (the)", "CG", "COG", "178"},
	{"Congo (the Democratic Republic of the)", "CD", "COD", "180"},
	{"Cook Islands (the)", "CK", "COK", "184"},
	{"Costa Rica", "CR", "CRI", "188"},
	{"Croatia", "HR", "HRV", "191"},
	{"Cuba", "CU", "CUB", "192"},
	{"Cyprus", "CY", "CYP", "196"},
	{"Czech Republic (the)", "CZ", "CZE", "203"},
	{"Benin", "BJ", "BEN", "204"},
	{"Denmark", "DK", "DNK", "208"},
	{"Dominica", "DM", "DMA", "212"},
	{"Dominican Republic (the)", "DO", "DOM", "214"},
	{"Ecuador", "EC", "ECU", "218"},
	{"El Salvador", "SV", "SLV", "222"},
	{"Equatorial Guinea", "GQ", "GNQ", "226"},
	{"Ethiopia", "ET", "ETH", "231"},
	{"Eritrea", "ER", "ERI", "232"},
	{"Estonia", "EE", "EST", "233"},
	{"Faroe Islands (the)", "FO", "FRO", "234"},
	{"Falkland Islands (the) [Malvinas]", "FK", "FLK", "238"},
	{"South Georgia and the South Sandwich Islands", "GS", "SGS", "239"},
	{"Fiji", "FJ", "FJI", "242"},
	{"Finland", "FI", "FIN", "246"},
	{"Åland Islands", "AX", "ALA", "248"},
	{"France", "FR", "FRA", "250"},
	{"French Guiana", "GF", "GUF", "254"},
	{"French Polynesia", "PF", "PYF", "258"},
	{"French Southern Territories (the)", "TF", "ATF", "260"},
	{"Djibouti", "DJ", "DJI", "262"},
	{"Gabon", "GA", "GAB", "266"},
	{"Georgia", "GE", "GEO", "268"},
	{"Gambia (the)", "GM", "GMB", "270"},
	{"Palestine, State of", "PS", "PSE", "275"},
	{"Germany", "DE", "DEU", "276"},
	{"Ghana", "GH", "GHA", "288"},
	{"Gibraltar", "GI", "GIB", "292"},
	{"Kiribati", "KI", "KIR", "296"},
	{"Greece", "GR", "GRC", "300"},
	{"Greenland", "GL", "GRL", "304"},
	{"Grenada", "GD", "GRD", "308"},
	{"Guadeloupe", "GP", "GLP", "312"},
	{"Guam", "GU", "GUM", "316"},
	{"Guatemala", "GT", "GTM", "320"},
	{"Guinea", "GN", "GIN", "324"},
	{"Guyana", "GY", "GUY", "328"},
	{"Haiti", "HT", "HTI", "332"},
	{"Heard Island and McDonald Islands", "HM", "HMD", "334"},
	{"Holy See (the)", "VA", "VAT", "336"},
	{"Honduras", "HN", "HND", "340"},
	{"Hong Kong", "HK", "HKG", "344"},
	{"Hungary", "HU", "HUN", "348"},
	{"Iceland", "IS", "ISL", "352"},
	{"India", "IN", "IND", "356"},
	{"Indonesia", "ID", "IDN", "360"},
	{"Iran (Islamic Republic of)", "IR", "IRN", "364"},
	{"Iraq", "IQ", "IRQ", "368"},
	{"Ireland", "IE", "IRL", "372"},
	{"Israel", "IL", "ISR", "376"},
	{"Italy", "IT", "ITA", "380"},
	{"Côte d'Ivoire", "CI", "CIV", "384"},
	{"Jamaica", "JM", "JAM", "388"},
	{"Japan", "JP", "JPN", "392"},
	{"Kazakhstan", "KZ", "KAZ", "398"},
	{"Jordan", "JO", "JOR", "400"},
	{"Kenya", "KE", "KEN", "404"},
	{"Korea (the Democratic People's Republic of)", "KP", "PRK", "408"},
	{"Korea (the Republic of)", "KR", "KOR", "410"},
	{"Kuwait", "KW", "KWT", "414"},
	{"Kyrgyzstan", "KG", "KGZ", "417"},
	{"Lao People's Democratic Republic (the)", "LA", "LAO", "418"},
	{"Lebanon", "LB", "LBN", "422"},
	{"Lesotho", "LS", "LSO", "426"},
	{"Latvia", "LV", "LVA", "428"},
	{"Liberia", "LR", "LBR", "430"},
	{"Libya", "LY", "LBY", "434"},
	{"Liechtenstein", "LI", "LIE", "438"},
	{"Lithuania", "LT", "LTU", "440"},
	{"Luxembourg", "LU", "LUX", "442"},
	{"Macao", "MO", "MAC", "446"},
	{"Madagascar", "MG", "MDG", "450"},
	{"Malawi", "MW", "MWI", "454"},
	{"Malaysia", "MY", "MYS", "458"},
	{"Maldives", "MV", "MDV", "462"},
	{"Mali", "ML", "MLI", "466"},
	{"Malta", "MT", "MLT", "470"},
	{"Martinique", "MQ", "MTQ", "474"},
	{"Mauritania", "MR", "MRT", "478"},
	{"Mauritius", "MU", "MUS", "480"},
	{"Mexico", "MX", "MEX", "484"},
	{"Monaco", "MC", "MCO", "492"},
	{"Mongolia", "MN", "MNG", "496"},
	{"Moldova (the Republic of)", "MD", "MDA", "498"},
	{"Montenegro", "ME", "MNE", "499"},
	{"Montserrat", "MS", "MSR", "500"},
	{"Morocco", "MA", "MAR", "504"},
	{"Mozambique", "MZ", "MOZ", "508"},
	{"Oman", "OM", "OMN", "512"},
	{"Namibia", "NA", "NAM", "516"},
	{"Nauru", "NR", "NRU", "520"},
	{"Nepal", "NP", "NPL", "524"},
	{"Netherlands (the)", "NL", "NLD", "528"},
	{"Curaçao", "CW", "CUW", "531"},
	{"Aruba", "AW", "ABW", "533"},
	{"Sint Maarten (Dutch part)", "SX", "SXM", "534"},
	{"Bonaire, Sint Eustatius and Saba", "BQ", "BES", "535"},
	{"New Caledonia", "NC", "NCL", "540"},
	{"Vanuatu", "VU", "VUT", "548"},
	{"New Zealand", "NZ", "NZL", "554"},
	{"Nicaragua", "NI", "NIC", "558"},
	{"Niger (the)", "NE", "NER", "562"},
	{"Nigeria", "NG", "NGA", "566"},
	{"Niue", "NU", "NIU", "570"},
	{"Norfolk Island", "NF", "NFK", "574"},
	{"Norway", "NO", "NOR", "578"},
	{"Northern Mariana Islands (the)", "MP", "MNP", "580"},
	{"United States Minor Outlying Islands (the)", "UM", "UMI", "581"},
	{"Micronesia (Federated States of)", "FM", "FSM", "583"},
	{"Marshall Islands (the)", "MH", "MHL", "584"},
	{"Palau", "PW", "PLW", "585"},
	{"Pakistan", "PK", "PAK", "586"},
	{"Panama", "PA", "PAN", "591"},
	{"Papua New Guinea", "PG", "PNG", "598"},
	{"Paraguay", "PY", "PRY", "600"},
	{"Peru", "PE", "PER", "604"},
	{"Philippines (the)", "PH", "PHL", "608"},
	{"Pitcairn", "PN", "PCN", "612"},
	{"Poland", "PL", "POL", "616"},
	{"Portugal", "PT", "PRT", "620"},
	{"Guinea-Bissau", "GW", "GNB", "624"},
	{"Timor-Leste", "TL", "TLS", "626"},
	{"Puerto Rico", "PR", "PRI", "630"},
	{"Qatar", "QA", "QAT", "634"},
	{"Réunion", "RE", "REU", "638"},
	{"Romania", "RO", "ROU", "642"},
	{"Russian Federation (the)", "RU", "RUS", "643"},
	{"Rwanda", "RW", "RWA", "646"},
	{"Saint Barthélemy", "BL", "BLM", "652"},
	{"Saint Helena, Ascension and Tristan da Cunha", "SH", "SHN", "654"},
	{"Saint Kitts and Nevis", "KN", "KNA", "659"},
	{"Anguilla", "AI", "AIA", "660"},
	{"Saint Lucia", "LC", "LCA", "662"},
	{"Saint Martin (French part)", "MF", "MAF", "663"},
	{"Saint Pierre and Miquelon", "PM", "SPM", "666"},
	{"Saint Vincent and the Grenadines", "VC", "VCT", "670"},
	{"San Marino", "SM", "SMR", "674"},
	{"Sao Tome and Principe", "ST", "STP", "678"},
	{"Saudi Arabia", "SA", "SAU", "682"},
	{"Senegal", "SN", "SEN", "686"},
	{"Serbia", "RS", "SRB", "688"},
	{"Seychelles", "SC", "SYC", "690"},
	{"Sierra Leone", "SL", "SLE", "694"},
	{"Singapore", "SG", "SGP", "702"},
	{"Slovakia", "SK", "SVK", "703"},
	{"Viet Nam", "VN", "VNM", "704"},
	{"Slovenia", "SI", "SVN", "705"},
	{"Somalia", "SO", "SOM", "706"},
	{"South Africa", "ZA", "ZAF", "710"},
	{"Zimbabwe", "ZW", "ZWE", "716"},
	{"Spain", "ES", "ESP", "724"},
	{"South Sudan", "SS", "SSD", "728"},
	{"Sudan (the)", "SD", "SDN", "729"},
	{"Western Sahara*", "EH", "ESH", "732"},
	{"Suriname", "SR", "SUR", "740"},
	{"Svalbard and Jan Mayen", "SJ", "SJM", "744"},
	{"Swaziland", "SZ", "SWZ", "748"},
	{"Sweden", "SE", "SWE", "752"},
	{"Switzerland", "CH", "CHE", "756"},
	{"Syrian Arab Republic", "SY", "SYR", "760"},
	{"Tajikistan", "TJ", "TJK", "762"},
	{"Thailand", "TH", "THA", "764"},
	{"Togo", "TG", "TGO", "768"},
	{"Tokelau", "TK", "TKL", "772"},
	{"Tonga", "TO", "TON", "776"},
	{"Trinidad and Tobago", "TT", "TTO", "780"},
	{"United Arab Emirates (the)", "AE", "ARE", "784"},
	{"Tunisia", "TN", "TUN", "788"},
	{"Turkey", "TR", "TUR", "792"},
	{"Turkmenistan", "TM", "TKM", "795"},
	{"Turks and Caicos Islands (the)", "TC", "TCA", "796"},
	{"Tuvalu", "TV", "TUV", "798"},
	{"Uganda", "UG", "UGA", "800"},
	{"Ukraine", "UA", "UKR", "804"},
	{"Macedonia (the former Yugoslav Republic of)", "MK", "MKD", "807"},
	{"Egypt", "EG", "EGY", "818"},
	{"United Kingdom of Great Britain and Northern Ireland (the)", "GB", "GBR", "826"},
	{"Guernsey", "GG", "GGY", "831"},
	{"Jersey", "JE", "JEY", "832"},
	{"Isle of Man", "IM", "IMN", "833"},
	{"Tanzania, United Republic of", "TZ", "TZA", "834"},
	{"United States of America (the)", "US", "USA", "840"},
	{"Virgin Islands (U.S.)", "VI", "VIR", "850"},
	{"Burkina Faso", "BF", "BFA", "854"},
	{"Uruguay", "UY", "URY", "858"},
	{"Uzbekistan", "UZ", "UZB", "860"},
	{"Venezuela (Bolivarian Republic of)", "VE", "VEN", "862"},
	{"Wallis and Futuna", "WF", "WLF", "876"},
	{"Samoa", "WS", "WSM", "882"},
	{"Yemen", "YE", "YEM", "887"},
	{"Zambia", "ZM", "ZMB", "894"},
}

func ISO3166Constraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	err := StringConstraint(ctx)
	if err != nil {
		return err
	}

	str := ctx.Value().(string)
	if str == "" {
		return nil
	}

	compare, _ := ctx.Arg().(string)
	for _, country := range iso3166List {
		switch compare {
		case "alpha3":
			if country.Alpha3 == ctx.Value() {
				return nil
			}
		case "num":
			if country.Numeric == ctx.Value() {
				return nil
			}
		default:
			if country.Alpha2 == ctx.Value() {
				return nil
			}
		}
	}

	ctx.BuildViolation("invalid_country", map[string]interface{}{
		"kind":  compare,
		"value": ctx.Value(),
	}).AddViolation()

	return nil
}

func ISO639Constraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	err := StringConstraint(ctx)
	if err != nil {
		return err
	}

	str := ctx.Value().(string)
	if str == "" {
		return nil
	}

	if _, err := language.Parse(ctx.Value().(string)); err != nil {
		ctx.BuildViolation("invalid_language", map[string]interface{}{
			"value": ctx.Value(),
		}).AddViolation()
	}

	return nil
}
