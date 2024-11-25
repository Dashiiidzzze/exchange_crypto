package app

// добавление лотов в таблицу lot и pair
func Init(pairList []string) {
	for i := 0; i < len(pairList); i++ {
		var reqBDcheck string = "SELECT * FROM lot WHERE lot.name = '" + pairList[i] + "'"
		response, err := RquestDataBase(reqBDcheck)
		if err != nil {
			return
		} else if response == "" {
			var reqBD string = "INSERT INTO lot VALUES ('" + pairList[i] + "')"

			_, err2 := RquestDataBase(reqBD)
			if err2 != nil {
				return
			}
		}
	}
	var allLots []string
	for i := 0; i < len(pairList); i++ {
		var reqBDcheck string = "SELECT lot.lot_id FROM lot WHERE lot.name = '" + pairList[i] + "'"
		response, err := RquestDataBase(reqBDcheck)
		if err != nil {
			return
		}
		response = response[:len(response)-2]
		allLots = append(allLots, response)
	}
	for i := 0; i < len(allLots); i++ {
		for j := i + 1; j < len(allLots); j++ {
			var reqBDcheck string = "SELECT * FROM pair WHERE pair.first_lot_id = '" + allLots[i] + "' AND pair.second_lot_id = '" + allLots[j] + "'"
			response, err := RquestDataBase(reqBDcheck)
			if err != nil {
				return
			} else if response == "" {
				var reqBD string = "INSERT INTO pair VALUES ('" + allLots[i] + "', '" + allLots[j] + "')"

				_, err2 := RquestDataBase(reqBD)
				if err2 != nil {
					return
				}
			}
		}
	}
}
