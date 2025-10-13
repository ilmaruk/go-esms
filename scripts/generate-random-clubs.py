import json
import random
import string

leagues = [
    {
        "name": "A",
        "groups": [""],
        "count": 20,
        "elo_max": 1929,
        "elo_min": 1531,
    },
    {
        "name": "B",
        "groups": ["A", "B"],
        "count": 16,
        "elo_max": 1597,
        "elo_min": 1228,
    },
    {
        "name": "C",
        "groups": ["A", "B", "C"],
        "count": 16,
        "elo_max": 1277,
        "elo_min": 1022,
    },
]

clubs = [
  {"name":"Spartak Velgorod","colors":{"primary":"#C8102E","secondary":"#FFFFFF","accent":"#0033A0"}},
  {"name":"Dynamo Kalenov","colors":{"primary":"#0047AB","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Lokomotiv Dravik","colors":{"primary":"#006400","secondary":"#FFFFFF","accent":"#FFCC00"}},
  {"name":"Slavia Norven","colors":{"primary":"#8B0000","secondary":"#F5F5F5","accent":"#000000"}},
  {"name":"Partizan Trevna","colors":{"primary":"#000000","secondary":"#FFFFFF","accent":"#B22222"}},
  {"name":"Metalist Dorvak","colors":{"primary":"#1E90FF","secondary":"#002244","accent":"#FFD166"}},
  {"name":"Torpedo Ruznar","colors":{"primary":"#2E8B57","secondary":"#FFFFFF","accent":"#8B4513"}},
  {"name":"Karpaty Valenik","colors":{"primary":"#006B3C","secondary":"#FFFFFF","accent":"#F2C649"}},
  {"name":"Slovan Brevik","colors":{"primary":"#0B5EA6","secondary":"#FFFFFF","accent":"#E03A3E"}},
  {"name":"CSKA Dornov","colors":{"primary":"#B22222","secondary":"#003366","accent":"#FFD700"}},
  {"name":"Zorya Mirgrad","colors":{"primary":"#0A2342","secondary":"#FFFFFF","accent":"#F4D35E"}},
  {"name":"Hajduk Straven","colors":{"primary":"#8B0000","secondary":"#F0EDEE","accent":"#FFD700"}},
  {"name":"Tatran Orlovec","colors":{"primary":"#228B22","secondary":"#FFFFFF","accent":"#2F4F4F"}},
  {"name":"Sloga Drunov","colors":{"primary":"#800000","secondary":"#FFDAB9","accent":"#000000"}},
  {"name":"Rudar Zelenik","colors":{"primary":"#2F4F2F","secondary":"#C0C0C0","accent":"#FFD700"}},
  {"name":"Obolon Karvich","colors":{"primary":"#2C9AB7","secondary":"#FFFFFF","accent":"#0B3D91"}},
  {"name":"Botev Trusna","colors":{"primary":"#0066CC","secondary":"#FFFFFF","accent":"#FFCC00"}},
  {"name":"Radnički Malgor","colors":{"primary":"#A52A2A","secondary":"#FFFFFF","accent":"#000000"}},
  {"name":"Arda Lysvek","colors":{"primary":"#004B87","secondary":"#FFFFFF","accent":"#F4A261"}},
  {"name":"Górnik Dovran","colors":{"primary":"#1B3A4B","secondary":"#FFFFFF","accent":"#B8860B"}},
  {"name":"Vardar Kravena","colors":{"primary":"#8A1538","secondary":"#FFD89B","accent":"#000000"}},
  {"name":"Olimpija Tarnovik","colors":{"primary":"#006B6B","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Rotor Belvik","colors":{"primary":"#0D47A1","secondary":"#FFFFFF","accent":"#E53935"}},
  {"name":"Energetik Narven","colors":{"primary":"#FF8C00","secondary":"#FFFFFF","accent":"#004D40"}},
  {"name":"Cherno More Dreznik","colors":{"primary":"#004C6D","secondary":"#FFFFFF","accent":"#00A86B"}},
  {"name":"Severna Kralov","colors":{"primary":"#28334A","secondary":"#BFC8D6","accent":"#D4AF37"}},
  {"name":"Volga Pradnik","colors":{"primary":"#003366","secondary":"#FFFFFF","accent":"#FFCC00"}},
  {"name":"Sloboda Orshen","colors":{"primary":"#8B4513","secondary":"#FFF8E1","accent":"#2E8B57"}},
  {"name":"Rubin Drezha","colors":{"primary":"#8A1538","secondary":"#F5F5F5","accent":"#FFD700"}},
  {"name":"Metalurg Varnov","colors":{"primary":"#4B5320","secondary":"#FFFFFF","accent":"#C0C0C0"}},
  {"name":"Torpedo Blazna","colors":{"primary":"#1F3A93","secondary":"#FFFFFF","accent":"#FF5A5F"}},
  {"name":"Polonia Mirkov","colors":{"primary":"#FFFFFF","secondary":"#DC143C","accent":"#00008B"}},
  {"name":"Spartak Nardev","colors":{"primary":"#C8102E","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Dynamo Zelgor","colors":{"primary":"#0033A0","secondary":"#FFFFFF","accent":"#F4D35E"}},
  {"name":"Lokomotiv Krusha","colors":{"primary":"#006400","secondary":"#FFDD00","accent":"#FFFFFF"}},
  {"name":"Partizan Drezal","colors":{"primary":"#000000","secondary":"#FFFFFF","accent":"#B22222"}},
  {"name":"Slavia Korvenik","colors":{"primary":"#B22222","secondary":"#FFFFFF","accent":"#000000"}},
  {"name":"Tatran Bredov","colors":{"primary":"#2E8B57","secondary":"#FFFFFF","accent":"#8B4513"}},
  {"name":"Sloga Tavrin","colors":{"primary":"#800000","secondary":"#FFDAB9","accent":"#000000"}},
  {"name":"Ruch Kroslav","colors":{"primary":"#0033A0","secondary":"#FFFFFF","accent":"#FCCC0A"}},
  {"name":"Zaglebie Vorna","colors":{"primary":"#0066CC","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Vojvodina Dralov","colors":{"primary":"#990000","secondary":"#FFFFFF","accent":"#003366"}},
  {"name":"Slovan Greven","colors":{"primary":"#0B5EA6","secondary":"#FFFFFF","accent":"#E03A3E"}},
  {"name":"Mladost Ornar","colors":{"primary":"#1E7327","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Rudar Zelgrad","colors":{"primary":"#2F4F2F","secondary":"#C0C0C0","accent":"#8B0000"}},
  {"name":"Obolon Krevik","colors":{"primary":"#2C9AB7","secondary":"#FFFFFF","accent":"#0B3D91"}},
  {"name":"Karpaty Vornes","colors":{"primary":"#006B3C","secondary":"#FFFFFF","accent":"#FFD166"}},
  {"name":"CSKA Brezhin","colors":{"primary":"#B22222","secondary":"#003366","accent":"#FFD700"}},
  {"name":"Torpedo Marnik","colors":{"primary":"#2E8B57","secondary":"#FFFFFF","accent":"#8B4513"}},
  {"name":"Spartak Traven","colors":{"primary":"#C8102E","secondary":"#FFFFFF","accent":"#000000"}},
  {"name":"Dynamo Lornik","colors":{"primary":"#0047AB","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Lokomotiv Draskar","colors":{"primary":"#006400","secondary":"#FFFFFF","accent":"#FFCC00"}},
  {"name":"Partizan Molgor","colors":{"primary":"#000000","secondary":"#FFFFFF","accent":"#B22222"}},
  {"name":"Slavia Tresnik","colors":{"primary":"#8B0000","secondary":"#F5F5F5","accent":"#000000"}},
  {"name":"Metalist Norvak","colors":{"primary":"#1E90FF","secondary":"#002244","accent":"#FFD166"}},
  {"name":"Radnički Zelven","colors":{"primary":"#A52A2A","secondary":"#FFFFFF","accent":"#000000"}},
  {"name":"Hajduk Vrasna","colors":{"primary":"#8B0000","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Sloga Durven","colors":{"primary":"#800000","secondary":"#FFDAB9","accent":"#2E8B57"}},
  {"name":"Zorya Blenik","colors":{"primary":"#0A2342","secondary":"#FFFFFF","accent":"#F4D35E"}},
  {"name":"Tatran Drunov","colors":{"primary":"#228B22","secondary":"#FFFFFF","accent":"#2F4F4F"}},
  {"name":"Vardar Molgrad","colors":{"primary":"#8A1538","secondary":"#FFD89B","accent":"#000000"}},
  {"name":"Volga Orsden","colors":{"primary":"#003366","secondary":"#FFFFFF","accent":"#FFCC00"}},
  {"name":"Rotor Kalven","colors":{"primary":"#0D47A1","secondary":"#FFFFFF","accent":"#E53935"}},
  {"name":"Energetik Tavros","colors":{"primary":"#FF8C00","secondary":"#FFFFFF","accent":"#004D40"}},
  {"name":"Sloboda Kriven","colors":{"primary":"#8B4513","secondary":"#FFF8E1","accent":"#2E8B57"}},
  {"name":"Spartak Dornik","colors":{"primary":"#C8102E","secondary":"#FFFFFF","accent":"#0033A0"}},
  {"name":"Dynamo Bresnar","colors":{"primary":"#0047AB","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Lokomotiv Varlen","colors":{"primary":"#006400","secondary":"#FFDD00","accent":"#FFFFFF"}},
  {"name":"Partizan Zelgor","colors":{"primary":"#000000","secondary":"#FFFFFF","accent":"#B22222"}},
  {"name":"Slavia Branik","colors":{"primary":"#B22222","secondary":"#FFFFFF","accent":"#000000"}},
  {"name":"Karpaty Norvash","colors":{"primary":"#006B3C","secondary":"#FFFFFF","accent":"#F2C649"}},
  {"name":"Metalurg Praven","colors":{"primary":"#4B5320","secondary":"#FFFFFF","accent":"#C0C0C0"}},
  {"name":"Hajduk Krezna","colors":{"primary":"#8B0000","secondary":"#F0EDEE","accent":"#FFD700"}},
  {"name":"Sloga Dornesk","colors":{"primary":"#800000","secondary":"#FFDAB9","accent":"#000000"}},
  {"name":"Rudar Velmar","colors":{"primary":"#2F4F2F","secondary":"#C0C0C0","accent":"#FFD700"}},
  {"name":"CSKA Broven","colors":{"primary":"#B22222","secondary":"#003366","accent":"#FFD700"}},
  {"name":"Torpedo Drasnik","colors":{"primary":"#1F3A93","secondary":"#FFFFFF","accent":"#FF5A5F"}},
  {"name":"Rubin Kovnar","colors":{"primary":"#8A1538","secondary":"#F5F5F5","accent":"#FFD700"}},
  {"name":"Volga Mirvos","colors":{"primary":"#003366","secondary":"#FFFFFF","accent":"#FFCC00"}},
  {"name":"Spartak Trevon","colors":{"primary":"#C8102E","secondary":"#FFFFFF","accent":"#000000"}},
  {"name":"Dynamo Dalvik","colors":{"primary":"#0047AB","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Lokomotiv Brezna","colors":{"primary":"#006400","secondary":"#FFFFFF","accent":"#FFCC00"}},
  {"name":"Partizan Korvik","colors":{"primary":"#000000","secondary":"#FFFFFF","accent":"#B22222"}},
  {"name":"Slavia Mornav","colors":{"primary":"#8B0000","secondary":"#F5F5F5","accent":"#000000"}},
  {"name":"Tatran Zelvik","colors":{"primary":"#228B22","secondary":"#FFFFFF","accent":"#2F4F4F"}},
  {"name":"Zorya Blornik","colors":{"primary":"#0A2342","secondary":"#FFFFFF","accent":"#F4D35E"}},
  {"name":"Sloga Pradven","colors":{"primary":"#800000","secondary":"#FFDAB9","accent":"#2E8B57"}},
  {"name":"Karpaty Drusk","colors":{"primary":"#006B3C","secondary":"#FFFFFF","accent":"#F2C649"}},
  {"name":"Metalist Tavrin","colors":{"primary":"#1E90FF","secondary":"#002244","accent":"#FFD166"}},
  {"name":"Vardar Kolven","colors":{"primary":"#8A1538","secondary":"#FFD89B","accent":"#000000"}},
  {"name":"Radnički Drosmar","colors":{"primary":"#A52A2A","secondary":"#FFFFFF","accent":"#000000"}},
  {"name":"Spartak Krelov","colors":{"primary":"#C8102E","secondary":"#FFFFFF","accent":"#0033A0"}},
  {"name":"Dynamo Tvarin","colors":{"primary":"#0047AB","secondary":"#FFFFFF","accent":"#FFD700"}},
  {"name":"Lokomotiv Ornash","colors":{"primary":"#006400","secondary":"#FFFFFF","accent":"#FFCC00"}},
  {"name":"Slavia Narsk","colors":{"primary":"#8B0000","secondary":"#F5F5F5","accent":"#000000"}},
  {"name":"Partizan Velgrad","colors":{"primary":"#000000","secondary":"#FFFFFF","accent":"#B22222"}},
  {"name":"CSKA Marnov","colors":{"primary":"#B22222","secondary":"#003366","accent":"#FFD700"}},
  {"name":"Rudar Treznik","colors":{"primary":"#2F4F2F","secondary":"#C0C0C0","accent":"#FFD700"}},
  {"name":"Torpedo Dovnar","colors":{"primary":"#1F3A93","secondary":"#FFFFFF","accent":"#FF5A5F"}},
  {"name":"Volga Kravich","colors":{"primary":"#003366","secondary":"#FFFFFF","accent":"#FFCC00"}}
]


def random_code(name: str, collision: bool = False) -> str:
    parts = name.split(" ")
    if len(parts) == 1:
        code = name[0:4] if not collision else f"{name[0:3]}{name[3:4]}"
    elif len(parts) == 2:
        code = f"{parts[0][0:2]}{parts[1][0:2]}" if not collision else f"{parts[0][0:2]}{parts[1][0:1]}{parts[1][2:3]}"
    elif len(parts) == 3:
        code = f"{parts[0][0:1]}{parts[1][0:1]}{parts[2][0:2]}" if not collision else f"{parts[0][0:1]}{parts[1][0:1]}{parts[2][0:1]}{parts[2][2:3]}"
    
    return code.upper()
    

def extract_city(name: str) -> str:
    return name.split(" ")[-1]

random.shuffle(clubs)

codes = {}
index = 0
for league in leagues:
    delta = (league["elo_max"] - league["elo_min"]) / (league["count"] - 1)
    for g in league["groups"]:
        for c in range(league["count"]):
            club = clubs[index]

            code = random_code(club["name"])
            while code in codes:
                print("collision", club["name"], code, codes[code])
                code = random_code(club["name"], True)
            codes[code] = club["name"]

            clubs[index]["code"] = code
            clubs[index]["city"] = extract_city(club["name"])
            # clubs[index]["elo"] = round(league["elo_max"] - delta * c),
            clubs[index]["elo"] = int(round((league["elo_max"] - delta * c) * (1.+(random.random()-.5)/100.)),)
            clubs[index]["league"] = league["name"] + g

            index += 1

with open("clubs.json", "wt") as fh:
    json.dump(clubs, fh, ensure_ascii=False)
