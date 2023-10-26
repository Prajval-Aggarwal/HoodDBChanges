package db

import (
	"fmt"

	"gorm.io/gorm"
)

func TriggerFunc(db *gorm.DB) {
	fmt.Println("Trtigger func called")
	triggerFunc := `CREATE OR REPLACE FUNCTION update_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;`

	err := db.Exec(triggerFunc).Error
	if err != nil {
		fmt.Println("error is", err)
	}
}

func TbaleTriggers(db *gorm.DB) {
	tablesList, _ := db.Migrator().GetTables()

	triggerFirstHalf := ` CREATE OR REPLACE TRIGGER set_updated_at
	BEFORE UPDATE ON 
	`

	triggerSecondHalf := ` FOR EACH ROW
	EXECUTE FUNCTION update_updated_at();`

	for _, table := range tablesList {
		trigger := triggerFirstHalf + table + triggerSecondHalf
		fmt.Println("Triger query is", trigger)
		err := db.Exec(trigger).Error
		if err != nil {
			fmt.Println("Error is:", err)
			return
		}
	}
}
