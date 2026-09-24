package fileio

// GetGeneralName returns a human-readable name for the general ID
func GetGeneralName(generalId uint16) (string, bool) {
	switch generalId {
	case 25171:
		return "Muhammad", true
	case 25172:
		return "Mahdi", true
	case 25180:
		return "Simo Hayha", true
	case 25182:
		return "Katukov", true
	case 25184:
		return "Anders", true
	case 25185:
		return "Juliusz", true
	case 25187:
		return "Kluge", true
	case 25188:
		return "Horton", true
	case 25191:
		return "Kimmel", true
	case 25193:
		return "Hermann Balck", true
	case 25199:
		return "Pavlov", true
	case 25200:
		return "Tolbukhin", true
	case 25201:
		return "Yeryomenko", true
	case 25202:
		return "Vasily", true
	case 25203:
		return "Weygand", true
	case 25204:
		return "Sikorski", true
	case 25205:
		return "Mikawa", true
	case 25206:
		return "Higuchi", true
	case 25207:
		return "Richthofen", true
	case 25208:
		return "Arnim", true
	case 25209:
		return "Schmidhuber", true
	case 25210:
		return "Manstein", true
	case 25211:
		return "Kuznetsov", true
	case 25212:
		return "Abrams", true
	case 25217:
		return "Leslie", true
	case 25224:
		return "Helfrich", true
	case 25226:
		return "Patton", true
	case 25259:
		return "Papagos", true
	case 25261:
		return "Lhagvasuren", true
	case 25263:
		return "Alberto", true
	case 25266:
		return "Nasser", true
	case 25317:
		return "Qasim", true
	case 29001:
		return "Zhang.Z.Z", true
	case 29002:
		return "Sun.L.R", true
	case 29003:
		return "Zhu.D", true
	case 29004:
		return "Peng.D.H", true
	case 29005:
		return "Slim", true
	case 29006:
		return "Mountbatten", true
	case 29007:
		return "Dowding", true
	case 29008:
		return "Alexander", true
	case 29009:
		return "Montgomery", true
	case 29010:
		return "Auchinleck", true
	case 29011:
		return "Messe", true
	case 29012:
		return "Badoglio", true
	case 29013:
		return "Graziani", true
	case 29014:
		return "Chuikov", true
	case 29015:
		return "Vasilevsky", true
	case 29016:
		return "Bagramyan", true
	case 29017:
		return "Kuznetsov", true
	case 29018:
		return "Zhukov", true
	case 29019:
		return "Konev", true
	case 29020:
		return "Govorov", true
	case 29021:
		return "Rokossovsky", true
	case 29022:
		return "Vatutin", true
	case 29023:
		return "Timoshenko", true
	case 29024:
		return "Yamashita", true
	case 29025:
		return "Kuribayashi", true
	case 29026:
		return "Nagumo", true
	case 29027:
		return "Yamamoto", true
	case 29028:
		return "Terauchi", true
	case 29029:
		return "Tito", true
	case 29030:
		return "MacArthur", true
	case 29031:
		return "Eisenhower", true
	case 29032:
		return "Nimitz", true
	case 29033:
		return "Halsey", true
	case 29034:
		return "Fletcher", true
	case 29035:
		return "Arnold", true
	case 29036:
		return "Bradley", true
	case 29037:
		return "Patton", true
	case 29038:
		return "Crerar", true
	case 29039:
		return "Mannerheim", true
	case 29040:
		return "Tassigny", true
	case 29041:
		return "de Gaulle", true
	case 29042:
		return "Leclerc", true
	case 29043:
		return "Bock", true
	case 29044:
		return "Heinrici", true
	case 29045:
		return "Rundstedt", true
	case 29046:
		return "Donitz", true
	case 29047:
		return "Kesselring", true
	case 29048:
		return "Leeb", true
	case 29049:
		return "Guderian", true
	case 29050:
		return "Manstein", true
	case 29051:
		return "Rommel", true
	case 29052:
		return "Model", true
	case 29053:
		return "Smigly", true
	case 29054:
		return "Blamey", true
	case 29055:
		return "Nasser", true
	case 29056:
		return "Goring", true
	case 29057:
		return "Brauchitsch", true
	case 29058:
		return "Raeder", true
	case 29059:
		return "List", true
	case 29060:
		return "Student", true
	case 29061:
		return "Schorner", true
	case 29062:
		return "Weichs", true
	case 29063:
		return "Kuchler", true
	case 29064:
		return "Kleist", true
	case 29065:
		return "Manteuffel", true
	case 29066:
		return "Keitel", true
	case 29067:
		return "Witzleben", true
	case 29068:
		return "Greim", true
	case 29069:
		return "Kluge", true
	case 29070:
		return "Paulus", true
	case 29071:
		return "Meyer", true
	case 29072:
		return "Petain", true
	case 29073:
		return "Gamelin", true
	case 29074:
		return "Darlan", true
	case 29075:
		return "Juin", true
	case 29076:
		return "Leopold", true
	case 29077:
		return "Winkelman", true
	case 29078:
		return "Christian", true
	case 29079:
		return "Olav", true
	case 29080:
		return "Malinovsky", true
	case 29081:
		return "Voronov", true
	case 29082:
		return "Meretskov", true
	case 29083:
		return "Novikov", true
	case 29084:
		return "Sokolovsky", true
	case 29085:
		return "Shaposhnikov", true
	case 29086:
		return "Isakov", true
	case 29087:
		return "Voroshilov", true
	case 29088:
		return "Budyonny", true
	case 29089:
		return "Antonescu", true
	case 29090:
		return "Dumitrescu", true
	case 29091:
		return "Horthy", true
	case 29092:
		return "Riccardi", true
	case 29093:
		return "Campioni", true
	case 29094:
		return "Cavallero", true
	case 29095:
		return "Balbo", true
	case 29096:
		return "Cunningham", true
	case 29097:
		return "Wavell", true
	case 29098:
		return "Pound", true
	case 29099:
		return "Wingate", true
	case 29100:
		return "Dill", true
	case 29101:
		return "Portal", true
	case 29102:
		return "Dempsey", true
	case 29103:
		return "Papagos", true
	case 29104:
		return "Boris", true
	case 29105:
		return "Nagano", true
	case 29106:
		return "Okamura", true
	case 29107:
		return "Itagaki", true
	case 29108:
		return "Ushijima", true
	case 29109:
		return "Hata", true
	case 29110:
		return "Ozawa", true
	case 29111:
		return "Umezu", true
	case 29112:
		return "Kondo", true
	case 29113:
		return "Koga", true
	case 29114:
		return "Koiso", true
	case 29115:
		return "Yamaguchi", true
	case 29116:
		return "Toyoda", true
	case 29117:
		return "Inoue", true
	case 29118:
		return "Li.Z.R", true
	case 29119:
		return "Bai.C.X", true
	case 29120:
		return "Xue.Y", true
	case 29121:
		return "Liang.X.C", true
	case 29122:
		return "Du.Y.M", true
	case 29123:
		return "Zhang.L.F", true
	case 29124:
		return "Chen.S.K", true
	case 29125:
		return "Lin.B", true
	case 29126:
		return "Xu.S.Y", true
	case 29127:
		return "Liu.B.C", true
	case 29128:
		return "Xu.X.Q", true
	case 29129:
		return "Phibun", true
	case 29130:
		return "Thimayya", true
	case 29131:
		return "Crace", true
	case 29132:
		return "Franco", true
	case 29133:
		return "Clark", true
	case 29134:
		return "Spruance", true
	case 29135:
		return "Krueger", true
	case 29136:
		return "Devers", true
	case 29137:
		return "Eaker", true
	case 29138:
		return "Doolittle", true
	case 29139:
		return "Smith", true
	case 29140:
		return "King", true
	case 29141:
		return "Stilwell", true
	case 29142:
		return "Chennault", true
	case 29143:
		return "Kinkaid", true
	case 29144:
		return "Mitscher", true
	case 29145:
		return "McAuliffe", true
	case 29146:
		return "Lee", true
	case 29147:
		return "McCain", true
	case 29148:
		return "Simonds", true
	case 29151:
		return "Inonu", true
	case 29152:
		return "Dutra", true
	case 29153:
		return "Camacho", true
	case 29154:
		return "Chung", true
	case 29155:
		return "Choe", true
	case 29156:
		return "Osborn", true
	case 29157:
		return "Williams", true
	case 29158:
		return "Coulson", true
	case 29159:
		return "Wagner", true
	case 29160:
		return "Yudintsev", true
	case 29161:
		return "Gaiman", true
	case 29162:
		return "Kotick", true
	case 29163:
		return "Guillemot", true
	case 29164:
		return "Maldini", true
	case 29165:
		return "Morita", true
	case 29166:
		return "Benteke", true
	case 29167:
		return "Wu.R", true
	case 29168:
		return "Owairan", true
	case 29169:
		return "Davis", true
	default:
		return "", false
	}
}
