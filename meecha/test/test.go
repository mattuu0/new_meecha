package main

import (
	"log"
	"meecha/database"
	"meecha/location"

	"github.com/redis/go-redis/v9"

	"math"

)

func googleGeosail(LatA, LngA, LatB, LngB float64) float64 {
	i := math.Pi / 180
	r := 6371.008
	X := math.Acos(math.Sin(LatA*i)*math.Sin(LatB*i)+math.Cos(LatA*i)*math.Cos(LatB*i)*math.Cos(LngA*i-LngB*i)) * r
	return X
}

func main() {
	database.DBpath = "./test.db"
	database.Init()
	location.Init("test")

	ignore_points,err := location.Get_Ignore_Points("test2")

	if err != nil {
		log.Println(err)
		return
	}

	var keys []string
	for key := range ignore_points {
		keys = append(keys, key)
	}

	location.Remove_Ignore_Point("test2",keys[0])

	point := location.Ignore_point{
		Latitude: 34.730792491372604,
		Longitude: 135.59393416756606,
		Distance: 18300,
	}

	pointid, err := location.Add_Ignore_Point("test2",point)

	if err != nil {
		log.Println(err)
		return
	}

	point2 := location.Ignore_point{
		Latitude: 64.730792491372604,
		Longitude: 138.59393416756606,
		Distance: 100,
	}

	location.Update_Ignore_Point("test2",pointid,point2)
	
	log.Println(location.Get_Ignore_Points("test2"))
	//location.Remove_Ignore_Point("test", pointid)
}

func GeoAdd(client *redis.Client, key string, geoLocation ...*redis.GeoLocation) error {

	return nil
}

// 地点を表す構造体
type Point struct {
	Lat float64 // 緯度（度数法）
	Lon float64 // 経度（度数法）
}

// ラジアンに変換する関数
// deg: 度数法で表された角度
func toRad(deg float64) float64 {
	return deg * math.Pi / 180 // 度数法からラジアンに変換
}

// ビンセントィ公式を用いて2地点間の距離を求める関数
// p1, p2: 距離を計算したい2地点
func vincenty(p1, p2 Point) float64 {
	const a = 6378137.0                                // WGS-84楕円体の長半径（メートル）
	const b = 6356752.314245                           // WGS-84楕円体の短半径（メートル）
	const f = 1 / 298.257223563                        // WGS-84楕円体の平坦化
	L := toRad(p2.Lon - p1.Lon)                        // 2地点間の経度差（ラジアン）
	U1 := math.Atan((1 - f) * math.Tan(toRad(p1.Lat))) // 補助球上の緯度（ラジアン）
	U2 := math.Atan((1 - f) * math.Tan(toRad(p2.Lat))) // 補助球上の緯度（ラジアン）
	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)
	lambda := L                // 中心角の初期値
	for i := 0; i < 100; i++ { // 収束するまでループ
		sinLambda := math.Sin(lambda)
		cosLambda := math.Cos(lambda)
		sinSigma := math.Sqrt((cosU2*sinLambda)*(cosU2*sinLambda) + (cosU1*sinU2-sinU1*cosU2*cosLambda)*(cosU1*sinU2-sinU1*cosU2*cosLambda)) // 大円弧の中心角の正弦
		cosSigma := sinU1*sinU2 + cosU1*cosU2*cosLambda                                                                                      // 大円弧の中心角の余弦
		sigma := math.Atan2(sinSigma, cosSigma)                                                                                              // 大円弧の中心角
		sinAlpha := cosU1 * cosU2 * sinLambda / sinSigma                                                                                     // 方位角の正弦
		cosSqAlpha := 1 - sinAlpha*sinAlpha                                                                                                  // 方位角の余弦の二乗
		cos2SigmaM := cosSigma - 2*sinU1*sinU2/cosSqAlpha                                                                                    // 大円弧の中心角の余弦の二乗
		C := f / 16 * cosSqAlpha * (4 + f*(4-3*cosSqAlpha))                                                                                  // 補正項
		lambdaPrev := lambda                                                                                                                 // 前回の中心角
		lambda = L + (1-C)*f*sinAlpha*(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))                                // 中心角の更新
		if math.Abs(lambda-lambdaPrev) < 1e-12 {                                                                                             // 収束判定
			uSq := cosSqAlpha * (a*a - b*b) / (b * b)                                                                                                                    // 補助関数
			A := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq)))                                                                                                       // 補助関数
			B := uSq / 1024 * (256 + uSq*(-128+uSq*(74-47*uSq)))                                                                                                         // 補助関数
			deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM))) // 補正項
			s := b * A * (sigma - deltaSigma)                                                                                                                            // 地点間の距離
			return s                                                                                                                                                     // m
		}
	}
	return 0 // 収束しない場合は0を返す
}

// ハバーサイン公式を用いて2地点間の距離を求める
func haversine(p1, p2 Point) float64 {
	const R = 6371 // 地球の半径(キロメートル)
	dLat := toRad(p2.Lat - p1.Lat)
	dLon := toRad(p2.Lon - p1.Lon)
	lat1 := toRad(p1.Lat)
	lat2 := toRad(p2.Lat)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return (R * c) * 1000
}
