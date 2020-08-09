package views

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body", 
		spago.Tag("header", 			
			spago.A("class", spago.S(`navbar`)),
			spago.A("style", spago.S(`box-shadow: lightgrey 2px 2px 2px; padding: 1rem;`)),
			spago.Tag("section", 				
				spago.A("class", spago.S(`navber-section`)),
				spago.Tag("a", 					
					spago.A("class", spago.S(`navbar-brand mr-2`)),
					spago.A("style", spago.S(`text-transform: uppercase; font-weight: bold;`)),
					spago.T(`MT-BAND Monitor`),
				),
			),
		),
		spago.Tag("main", 			
			spago.A("class", spago.S(`container`)),
			spago.A("style", spago.S(`padding: 1rem;`)),
			spago.Tag("form", 				
				spago.A("class", spago.S(`form-horizontal`)),
				spago.Tag("div", 					
					spago.A("class", spago.S(`form-group`)),
					spago.Tag("div", 						
						spago.A("class", spago.S(`col-3 col-sm-12`)),
						spago.Tag("label", 							
							spago.A("class", spago.S(`form-label`)),
							spago.A("for", spago.S(`input-example-1`)),
							spago.T(``, spago.S(spago.T("Device:", c.BLE.FwRevision)), ``),
						),
					),
					spago.Tag("div", 						
						spago.A("class", spago.S(`col-9 col-sm-12`)),
						spago.Tag("div", 							
							spago.A("class", spago.S(`btn-group`)),
							spago.Tag("button", 								
								spago.Event("click", c.OnStartClick),
								spago.A("class", spago.S(`btn btn-primary`)),
								spago.ClassMap{"disabled":c.BLE.IsConnect()},
								spago.T(`Connect`),
							),
							spago.Tag("button", 								
								spago.Event("click", c.OnStopClick),
								spago.A("class", spago.S(`btn`)),
								spago.ClassMap{"disabled":!c.BLE.IsConnect()},
								spago.T(`Disconnect`),
							),
						),
					),
				),
				spago.Tag("div", 					
					spago.A("class", spago.S(`form-group`)),
					spago.Tag("div", 						
						spago.A("class", spago.S(`col-3 col-sm-12`)),
						spago.Tag("label", 							
							spago.A("class", spago.S(`form-label`)),
							spago.A("for", spago.S(`input-example-1`)),
							spago.T(`Battery Remain`),
						),
					),
					spago.Tag("div", 						
						spago.A("class", spago.S(`col-9 col-sm-12`)),
						spago.A("style", spago.S(`margin: auto;`)),
						spago.Tag("div", 							
							spago.A("class", spago.S(`bar`)),
							spago.Tag("div", 								
								spago.A("class", spago.S(`bar-item`)),
								spago.A("style", spago.S(`width: `, spago.S(c.BLE.CurrentEnv.Battery), `%;`)),
								spago.T(``, spago.S(c.BLE.CurrentEnv.Battery), `%`),
							),
						),
					),
				),
				spago.Tag("div", 					
					spago.A("class", spago.S(`form-group`)),
					spago.Tag("div", 						
						spago.A("class", spago.S(`col-3 col-sm-12`)),
						spago.Tag("label", 							
							spago.A("class", spago.S(`form-label`)),
							spago.A("for", spago.S(`input-example-1`)),
							spago.T(`Sensor Values`),
						),
					),
					spago.Tag("div", 						
						spago.A("class", spago.S(`col-9 col-sm-12`)),
						spago.Tag("div", 							
							spago.A("class", spago.S(`columns`)),
							spago.Tag("div", 								
								spago.A("class", spago.S(`input-group column col-4 col-xs-12`)),
								spago.Tag("span", 									
									spago.A("class", spago.S(`input-group-addon addon-lg`)),
									spago.T(`BPM`),
								),
								
								spago.Tag("input", 									
									spago.A("type", spago.S(`text`)),
									spago.A("class", spago.S(`form-input input-lg`)),
									spago.A("readonly", spago.S(`true`)),
									spago.A("value", spago.S(``, spago.S(c.BLE.BPM), ``)),
								),
							),
							spago.Tag("div", 								
								spago.A("class", spago.S(`input-group column col-4 col-xs-12`)),
								spago.Tag("span", 									
									spago.A("class", spago.S(`input-group-addon addon-lg`)),
									spago.T(`Humidity`),
								),
								
								spago.Tag("input", 									
									spago.A("type", spago.S(`text`)),
									spago.A("class", spago.S(`form-input input-lg`)),
									spago.A("readonly", spago.S(`true`)),
									spago.A("value", spago.S(``, spago.S(c.BLE.CurrentEnv.GetHumidity()), ``)),
								),
							),
							spago.Tag("div", 								
								spago.A("class", spago.S(`input-group column col-4 col-xs-12`)),
								spago.Tag("span", 									
									spago.A("class", spago.S(`input-group-addon addon-lg`)),
									spago.T(`Temperature`),
								),
								
								spago.Tag("input", 									
									spago.A("type", spago.S(`text`)),
									spago.A("class", spago.S(`form-input input-lg`)),
									spago.A("readonly", spago.S(`true`)),
									spago.A("value", spago.S(``, spago.S(c.BLE.CurrentEnv.GetTemperature()), ``)),
								),
							),
							spago.Tag("div", 								
								spago.A("class", spago.S(`input-group column col-4 col-xs-12`)),
								spago.Tag("span", 									
									spago.A("class", spago.S(`input-group-addon addon-lg`)),
									spago.T(`Skin-Temp.`),
								),
								
								spago.Tag("input", 									
									spago.A("type", spago.S(`text`)),
									spago.A("class", spago.S(`form-input input-lg`)),
									spago.A("readonly", spago.S(`true`)),
									spago.A("value", spago.S(``, spago.S(c.BLE.CurrentEnv.GetSkinTemp()), ``)),
								),
							),
							spago.Tag("div", 								
								spago.A("class", spago.S(`input-group column col-4 col-xs-12`)),
								spago.Tag("span", 									
									spago.A("class", spago.S(`input-group-addon addon-lg`)),
									spago.T(`EST-Temp.`),
								),
								
								spago.Tag("input", 									
									spago.A("type", spago.S(`text`)),
									spago.A("class", spago.S(`form-input input-lg`)),
									spago.A("readonly", spago.S(`true`)),
									spago.A("value", spago.S(``, spago.S(c.BLE.CurrentEnv.GetEstTemp()), ``)),
								),
							),
							spago.Tag("div", 								
								spago.A("class", spago.S(`input-group column col-4 col-xs-12`)),
								spago.Tag("span", 									
									spago.A("class", spago.S(`input-group-addon addon-lg`)),
									spago.T(`Flags`),
								),
								
								spago.Tag("input", 									
									spago.A("type", spago.S(`text`)),
									spago.A("class", spago.S(`form-input input-lg`)),
									spago.A("readonly", spago.S(`true`)),
									spago.A("value", spago.S(``, spago.S(c.BLE.CurrentEnv.GetFlags()), ``)),
								),
							),
						),
					),
				),
			),
		),
	)
}
