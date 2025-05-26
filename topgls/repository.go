package topgls

import (
	"context"
	"hpc-express-service/utils"
	"time"

	"github.com/go-pg/pg/v9"
)

type Repository interface {
	// InsertManifestData(ctx context.Context, uploadLogUUID string, data []*UploadManifestModel, chunkSize int) error
	// GetAllManifestToPreImport(ctx context.Context, uploadLoggingUUID string) ([]*utils.GetDataManifestImport, error)
	InsertPreImportManifest(ctx context.Context, manifest *utils.InsertPreImportHeaderManifestModel, chunkSize int) error
	// GetOtherDataPreImport(ctx context.Context, goods, countryCode, shipperName, mawb string) (*GetOtherDataPreImportModel, error)
	GetMawb(ctx context.Context, timestamp string) (*utils.GetMawb, error)
	GetShipperBrands(ctx context.Context) ([]*GetShipperBrandModel, error)
	GetMasterHsCode(ctx context.Context) ([]*GetMasterHsCodeModel, error)
	GetFreightData(ctx context.Context, uploadLogUUID string) (*GetFreightDataModel, error)
}

type repository struct {
	contextTimeout time.Duration
}

func NewRepository(
	timeout time.Duration,
) Repository {
	return &repository{
		contextTimeout: timeout,
	}
}

// func (r repository) InsertManifestData(ctx context.Context, uploadLogUUID string, data []*UploadManifestModel, chunkSize int) error {
// 	db := ctx.Value("postgreSQLConn").(*pg.DB)
// 	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

// 	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	defer tx.Rollback()

// 	// Insert Manifest
// 	// Chunk slice
// 	chunked := utils.ChunkSlice(data, chunkSize)
// 	{

// 		for _, chunkedRows := range chunked {
// 			sqlStr := "INSERT INTO public.tbl_topgls_manifests (mawb, bag_no, hawb, hs_code, origin, shipper_name, consignee_name, weight, weight_unit, packaging, shipper_address, consignee_address, province, district, postcode, pcs, qty, goods, goods_en, currency, total_price, fob, freight, insurance, cif, category, duty, vat, cost, local_tracking, reference1, reference2, customer_code, shipper_tel, consignee_tel, dimension, dimension_repacking, width, length, height, volume_weight, timestamp, upload_logging_uuid) VALUES "
// 			vals := []interface{}{}
// 			for _, row := range chunkedRows {
// 				sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?::float8, ?, ?, ?, ?, ?, ?, ?, ?, ?::integer, ?, ?, ?, ?, ?::float8, ?::float8, ?::float8, ?::float8, ?, ?, ?::float8, ?::float8, ?, ?, ?, ?, ?, ?, ?, ?, ?::float8, ?::float8, ?::float8, ?::float8, ?, ?),"
// 				vals = append(vals, row.Mawb, utils.NewNullString(row.BagNo), utils.NewNullString(row.Hawb), utils.NewNullString(row.HsCode), utils.NewNullString(row.Origin), utils.NewNullString(row.ShipperName), utils.NewNullString(row.ConsigneeName), row.WgtValue, row.WgtUnit, utils.NewNullString(row.Packaging), utils.NewNullString(row.ShipperAddress), utils.NewNullString(row.ConsigneeAddress), utils.NewNullString(row.Province), utils.NewNullString(row.District), utils.NewNullString(row.Postcode), utils.NewNullString(row.Pcs), row.Qty, utils.NewNullString(row.Goods), utils.NewNullString(row.GoodsEN), row.Currency, row.TotalPrice, row.FOB, row.Freight, row.Insurance, row.CIF, utils.NewNullString(row.Cat), row.Duty, row.Vat, row.Cost, utils.NewNullString(row.LocalTrackingNo), utils.NewNullString(row.Reference1), utils.NewNullString(row.Reference2), utils.NewNullString(row.CustomerCode), utils.NewNullString(row.ShipperTel), utils.NewNullString(row.ConsigneeTel), utils.NewNullString(row.Dimension), utils.NewNullString(row.DimensionRepacking), row.Width, row.Length, row.Height, row.VolumeWeight, timestamp, uploadLogUUID)
// 			}

// 			// remove last comma,
// 			sqlStr = sqlStr[0 : len(sqlStr)-1]

// 			// Convert symbol ? to $
// 			sqlStr = utils.ReplaceSQL(sqlStr, "?")
// 			// sqlStr += " ON CONFLICT (local_no) DO NOTHING returning uuid, local_no;"

// 			// Prepare statement
// 			stmt, err := tx.Prepare(sqlStr)
// 			if err != nil {
// 				tx.Rollback()
// 				return err
// 			}
// 			defer stmt.Close()

// 			_, err = stmt.ExecContext(ctx, vals...)
// 			if err != nil {
// 				tx.Rollback()
// 				return err
// 			}
// 		}
// 	}

// 	tx.Commit()

// 	return nil
// }

// func (r repository) GetAllManifestToPreImport(ctx context.Context, uploadLoggingUUID string) ([]*utils.GetDataManifestImport, error) {
// 	db := ctx.Value("postgreSQLConn").(*pg.DB)
// 	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

// 	list := []*utils.GetDataManifestImport{}
// 	q := `
// 		WITH exchange_rate_cte AS (
// 				SELECT
// 						CASE
// 								WHEN mct."type" = 'inbound' THEN import_exchange_rate::numeric
// 								WHEN mct."type" = 'outbound' THEN export_exchange_rate::numeric
// 								ELSE 0
// 						end as exchange_rate
// 				FROM
// 						ship2cu.customs_exchange_rate
// 						join tbl_upload_loggings tul on tul.uuid = $1
// 						join master_convert_templates mct on mct.code  = tul.template_code
// 				WHERE
// 						currency_code = mct.currency
// 				LIMIT 1
// 		),
// 		freight_rate_cte AS (
// 				SELECT
// 						value::int AS freight_rate
// 				FROM
// 						configurations
// 				WHERE
// 						name = 'freight_rate'
// 				LIMIT 1
// 		)
// 		SELECT DISTINCT
// 				x.mawb,
// 				x.hawb,
// 				x.category,
// 				'' AS consignee_tax,
// 				'' AS consignee_branch,
// 				x.consignee_name AS consignee_name,
// 				x.consignee_address AS consignee_address,
// 				x.district AS consignee_district,
// 				'' AS consignee_sub_province,
// 				x.province AS consignee_province,
// 				x.postcode AS consignee_postcode,
// 				'TH' AS consignee_country_code,
// 				'' AS consignee_email,
// 				x.consignee_tel AS consignee_phone_number,
// 				x.shipper_name,
// 				x.shipper_address,
// 				sb.district AS shipper_district,
// 				sb.sub_district AS shipper_sub_province,
// 				sb.province AS shipper_province,
// 				sb.postal_code AS shipper_postcode,
// 				'KR' AS shipper_country_code,
// 				'' AS shipper_email,
// 				x.shipper_tel AS shipper_phone_number,
// 				mhc.hs_code AS tariff_code,
// 				mhc.tariff AS tariff_sequence,
// 				mhc.stat AS stat_code,
// 				x.goods AS en_description_goods,
// 				x.goods AS th_description_goods,
// 				x.qty AS qty,
// 				CASE
// 						WHEN mhc.unit_code = 'KGM' THEN 'C62'
// 						ELSE mhc.unit_code
// 				END AS qty_unit,
// 				x.weight AS net_weight,
// 				'KGM' AS net_weight_unit,
// 				x.weight AS gross_weight,
// 				'KGM' AS gross_weight_unit,
// 				'1' AS package,
// 				'PK' AS package_unit_code,
// 				(
// 						x.total_price * exchange_rate_cte.exchange_rate
// 				) + (
// 						x.weight * freight_rate_cte.freight_rate
// 				) + (
// 						(x.total_price * exchange_rate_cte.exchange_rate) * 0.01
// 				) AS cif_value_foreign,
// 				x.total_price AS fob_value_foreign,
// 				exchange_rate_cte.exchange_rate AS exchange_rate,
// 				x.mawb AS shipping_mark,
// 				'KR' AS consignment_country,
// 				x.weight * freight_rate_cte.freight_rate AS freight_value_foreign,
// 				'THB' AS freight_currency_code,
// 				(x.total_price * exchange_rate_cte.exchange_rate) * 0.01 AS insurance_value_foreign,
// 				'THB' AS insurance_currency_code,
// 				'' AS other_charge_value_foreign,
// 				'' AS other_charge_currency_code,
// 				x.hawb AS invoice_no,
// 				to_char(now() AT TIME ZONE 'utc' AT TIME ZONE 'Asia/Bangkok', 'DD/MM/YYYY') AS invoice_date
// 		FROM
// 				public.tbl_topgls_manifests x
// 		LEFT JOIN
// 				master_hs_code mhc ON mhc.goods_en = x.goods
// 		LEFT JOIN
// 				ship2cu.master_shipper_brands sb ON sb."name" = x.shipper_name
// 		CROSS JOIN
// 				exchange_rate_cte
// 		CROSS JOIN
// 				freight_rate_cte
// 		WHERE
// 				x.upload_logging_uuid = $1
// 		`

// 	// if offset != 0 && limit != 0 {
// 	// 	q += "LIMIT " + fmt.Sprintf("%v", limit) + " OFFSET " + fmt.Sprintf("%v", offset)
// 	// }

// 	stmt, err := db.Prepare(q)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.QueryContext(ctx, &list, uploadLoggingUUID)
// 	if err != nil {
// 		return list, err
// 	}

// 	return list, nil
// }

func (r repository) InsertPreImportManifest(ctx context.Context, manifest *utils.InsertPreImportHeaderManifestModel, chunkSize int) error {
	db := ctx.Value("postgreSQLConn").(*pg.DB)
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert HeaderManifest
	var headerUUID string
	{
		sqlStr :=
			`
			INSERT INTO public.tbl_pre_import_manifest_headers
				(
					upload_logging_uuid, discharge_port, vassel_name, arrival_date, customer_name
				)
			VALUES
				(
					?, ?, ?, ?, ?
				)
			RETURNING uuid
		`

		// Prepare statement
		stmt, err := tx.Prepare(utils.PrepareSQL(sqlStr))
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		values := []interface{}{}
		values = append(
			values,
			manifest.UploadLoggingUUID,
			manifest.DischargePort,
			manifest.VasselName,
			manifest.ArrivalDate,
			manifest.CustomerName,
		)

		_, err = stmt.QueryOneContext(ctx, &headerUUID, values...)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert Manifest
	// Chunk slice
	chunked := utils.ChunkSlice(manifest.Details, chunkSize)
	{

		for _, chunkedRows := range chunked {
			sqlStr :=
				`
				INSERT INTO public.tbl_pre_import_manifest_details 
					(
						header_uuid, master_air_waybill, house_air_waybill, category, consignee_tax, consignee_branch, consignee_name, consignee_address, consignee_district, consignee_subprovince, consignee_province, consignee_postcode, consignee_country_code, consignee_email, consignee_phone_number, shipper_name, shipper_address, shipper_district, shipper_subprovince, shipper_province, shipper_postcode, shipper_country_code, shipper_email, shipper_phone_number, tariff_code, tariff_sequence, statistical_code, english_description_of_good, thai_description_of_good, quantity, quantity_unit_code, net_weight, net_weight_unit_code, gross_weight, gross_weight_unit_code, package, package_unit_code, cif_value_foreign, fob_value_foreign, exchange_rate, currency_code, shipping_mark, consignment_country, freight_value_foreign, freight_currency_code, insurance_value_foreign, insurance_currency_code, other_charge_value_foreign, other_charge_currency_code, invoice_no, invoice_date
					) 
					VALUES 
			`
			vals := []interface{}{}
			for _, row := range chunkedRows {
				row.HeaderUUID = headerUUID

				sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
				vals = append(vals,
					utils.NewNullString(row.HeaderUUID),
					utils.NewNullString(row.MasterAirWaybill),
					utils.NewNullString(row.HouseAirWaybill),
					utils.NewNullString(row.Category),
					utils.NewNullString(row.ConsigneeTax),
					utils.NewNullString(row.ConsigneeBranch),
					utils.NewNullString(row.ConsigneeName),
					utils.NewNullString(row.ConsigneeAddress),
					utils.NewNullString(row.ConsigneeDistrict),
					utils.NewNullString(row.ConsigneeSubprovince),
					utils.NewNullString(row.ConsigneeProvince),
					utils.NewNullString(row.ConsigneePostcode),
					utils.NewNullString(row.ConsigneeCountryCode),
					utils.NewNullString(row.ConsigneeEmail),
					utils.NewNullString(row.ConsigneePhoneNumber),
					utils.NewNullString(row.ShipperName),
					utils.NewNullString(row.ShipperAddress),
					utils.NewNullString(row.ShipperDistrict),
					utils.NewNullString(row.ShipperSubprovince),
					utils.NewNullString(row.ShipperProvince),
					utils.NewNullString(row.ShipperPostcode),
					utils.NewNullString(row.ShipperCountryCode),
					utils.NewNullString(row.ShipperEmail),
					utils.NewNullString(row.ShipperPhoneNumber),
					utils.NewNullString(row.TariffCode),
					utils.NewNullString(row.TariffSequence),
					utils.NewNullString(row.StatisticalCode),
					utils.NewNullString(row.EnglishDescriptionOfGood),
					utils.NewNullString(row.ThaiDescriptionOfGood),
					row.Quantity,
					utils.NewNullString(row.QuantityUnitCode),
					row.NetWeight,
					utils.NewNullString(row.NetWeightUnitCode),
					row.GrossWeight,
					utils.NewNullString(row.GrossWeightUnitCode),
					utils.NewNullString(row.Package),
					utils.NewNullString(row.PackageUnitCode),
					row.CifValueForeign,
					row.FobValueForeign,
					row.ExchangeRate,
					utils.NewNullString(row.CurrencyCode),
					utils.NewNullString(row.ShippingMark),
					utils.NewNullString(row.ConsignmentCountry),
					row.FreightValueForeign,
					utils.NewNullString(row.FreightCurrencyCode),
					row.InsuranceValueForeign,
					utils.NewNullString(row.InsuranceCurrencyCode),
					utils.NewNullString(row.OtherChargeValueForeign),
					utils.NewNullString(row.OtherChargeCurrencyCode),
					utils.NewNullString(row.InvoiceNo),
					utils.NewNullString(row.InvoiceDate),
				)
			}

			// remove last comma,
			sqlStr = sqlStr[0 : len(sqlStr)-1]

			// Convert symbol ? to $
			sqlStr = utils.ReplaceSQL(sqlStr, "?")
			// sqlStr += " ON CONFLICT (local_no) DO NOTHING returning uuid, local_no;"

			// Prepare statement
			stmt, err := tx.Prepare(sqlStr)
			if err != nil {
				tx.Rollback()
				return err
			}
			defer stmt.Close()

			_, err = stmt.ExecContext(ctx, vals...)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil
}

// func (r repository) GetOtherDataPreImport(ctx context.Context, goods, countryCode, shipperName, mawb string) (*GetOtherDataPreImportModel, error) {
// 	db := ctx.Value("postgreSQLConn").(*pg.DB)
// 	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

// 	x := GetOtherDataPreImportModel{}
// 	_, err := db.QueryOneContext(ctx, pg.Scan(
// 		&x.ShipperAddress,
// 		&x.ShipperDistrict,
// 		&x.ShipperSubprovince,
// 		&x.ShipperProvince,
// 		&x.ShipperPostcode,
// 		&x.TariffCode,
// 		&x.TariffSequence,
// 		&x.StatisticalCode,
// 		&x.QuantityUnitCode,
// 		&x.InvoiceDate,
// 	), `
// 		SELECT
// 			shipper.shipper_address,
// 			shipper.shipper_district,
// 			shipper.shipper_sub_province,
// 			shipper.shipper_province,
// 			shipper.shipper_postcode,
// 			mhsc.hs_code AS tariff_code,
// 			mhsc.tariff AS tariff_sequence,
// 			mhsc.stat AS statistical_code,
// 			CASE
// 					WHEN mhsc.unit_code = 'KGM' THEN 'C62'
// 					ELSE mhsc.unit_code
// 			END AS quantity_unit_code,
// 			TO_CHAR(NOW() AT TIME ZONE 'utc' AT TIME ZONE 'Asia/Bangkok', 'DD/MM/YYYY') AS invoice_date
// 		FROM
// 			(
// 				SELECT
// 					sb.address AS shipper_address,
// 					sb.district AS shipper_district,
// 					sb.sub_district AS shipper_sub_province,
// 					sb.province AS shipper_province,
// 					sb.postal_code AS shipper_postcode
// 				FROM ship2cu.master_shipper_brands sb
// 				WHERE sb.name = ?2
// 				AND sb.country_code = ?1
// 			) shipper

// 			NATURAL FULL JOIN

// 			(
// 				SELECT * FROM master_hs_code mhc WHERE mhc.goods_en = ?0
// 			) mhsc

// 	`, goods, countryCode, shipperName, mawb)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &x, nil
// }

func (r repository) GetShipperBrands(ctx context.Context) ([]*GetShipperBrandModel, error) {
	db := ctx.Value("postgreSQLConn").(*pg.DB)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	sqlStr := `
		SELECT
			sb.name AS shipper_name,
			sb.address AS shipper_address,
			sb.district AS shipper_district,
			sb.sub_district AS shipper_subprovince,
			sb.province AS shipper_province,
			sb.postal_code AS shipper_postcode,
			sb.country_code AS shipper_country_code
		FROM ship2cu.master_shipper_brands sb 
	`

	var list []*GetShipperBrandModel
	_, err := db.QueryContext(ctx, &list, sqlStr)

	if err != nil {
		return list, err
	}

	return list, nil

}

func (r repository) GetMasterHsCode(ctx context.Context) ([]*GetMasterHsCodeModel, error) {
	db := ctx.Value("postgreSQLConn").(*pg.DB)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	sqlStr := `
		SELECT 
			mhc.goods_en AS goods_en,
			mhc.hs_code AS tariff_code,
			mhc.tariff AS tariff_sequence,
			mhc.stat AS statistical_code,
			CASE
					WHEN mhc.unit_code = 'KGM' THEN 'C62'
					ELSE mhc.unit_code
			END AS quantity_unit_code
		FROM master_hs_code mhc 
	`

	var list []*GetMasterHsCodeModel
	_, err := db.QueryContext(ctx, &list, sqlStr)

	if err != nil {
		return list, err
	}

	return list, nil

}

func (r repository) GetMawb(ctx context.Context, mawb string) (*utils.GetMawb, error) {
	db := ctx.Value("postgreSQLConn").(*pg.DB)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	x := utils.GetMawb{}

	_, err := db.QueryOneContext(ctx, pg.Scan(
		&x.UUID,
		&x.FlightNo,
		&x.Origin,
		&x.Destination,
		&x.Mawb,
		&x.LotNo,
		&x.DepartureDatetime,
		&x.ArrivalDatetime,
		&x.Origin,
	), `
			SELECT 
				maw.uuid,
				maw.flight_no ,
				maw.origin_code ,
				maw.destination_code ,
				maw.lot_no_code ,
				maw.mawb,
				maw.lot_no_code,
				TO_CHAR(maw.departure_date_time at time zone 'utc' at time zone 'Asia/Bangkok', 'YYYY-MM-DD HH24:MI:SS') as departure_datetime,
				TO_CHAR(maw.arrival_date_time at time zone 'utc' at time zone 'Asia/Bangkok', 'YYYY-MM-DD HH24:MI:SS') as arrival_datetime
			FROM ship2cu.company_master_airway_bill maw 
			WHERE maw.mawb = ?
			AND maw.deleted_at IS NULL
			LIMIT 1
	 `, mawb)

	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (r repository) GetFreightData(ctx context.Context, uploadLogUUID string) (*GetFreightDataModel, error) {
	db := ctx.Value("postgreSQLConn").(*pg.DB)
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	x := GetFreightDataModel{}
	_, err := db.QueryOneContext(ctx, pg.Scan(
		&x.FreightZone,
		&x.FreightRate,
	), `
			SELECT (
				SELECT value FROM configurations c WHERE c.name = 'freight_zone'
			) AS freight_zone,
			(
				SELECT
						CASE
								WHEN mct."type" = 'inbound' THEN import_exchange_rate::numeric
								WHEN mct."type" = 'outbound' THEN export_exchange_rate::numeric
								ELSE 0
						end as exchange_rate
				FROM
						ship2cu.customs_exchange_rate
						join tbl_upload_loggings tul on tul.uuid = ?
						join master_convert_templates mct on mct.code  = tul.template_code
				WHERE
						currency_code = mct.currency
				LIMIT 1
			) AS freight_rate
	`, uploadLogUUID)

	if err != nil {
		return nil, err
	}

	return &x, nil
}
