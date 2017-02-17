package main

type item struct {
	id       string
	name     string
	itemType string
	x, y     string
}
type mapyCzLocation struct {
	rgeocode []item
}

const (
	xml string = `
<rgeocode label="část obce Karolinka, Karolinka, okres Vsetín" status="200" message="Ok">
<item id="4349" name="Karolinka" type="ward" x="18.2400622383" y="49.3512769199"/>
<item id="528" name="Karolinka" type="muni" x="18.2400622383" y="49.3512769199"/>
<item id="45" name="Vsetín" type="dist" x="18.0783387755" y="49.371466666"/>
<item id="9" name="Zlínský" type="regi" x="17.747371" y="49.220465"/>
<item id="112" name="Česko" type="coun" x="15.338411" y="49.742858"/>
</rgeocode>`
	link string = `https://api.mapy.cz/rgeocode?lat=49.3920400&lon=18.2485131`
)

func RgeocodeTest() {

}
