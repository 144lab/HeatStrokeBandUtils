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
					spago.T(`MT-BAND Dump Tool`),
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
								spago.T(`Dump`),
							),
							spago.Tag("button", 								
								spago.Event("click", c.OnStopClick),
								spago.A("class", spago.S(`btn`)),
								spago.ClassMap{"disabled":!c.BLE.IsConnect()},
								spago.T(`Abort`),
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
							spago.T(`Progress`),
						),
					),
					spago.Tag("div", 						
						spago.A("class", spago.S(`col-9 col-sm-12`)),
						spago.A("style", spago.S(`margin: auto;`)),
						spago.Tag("div", 							
							spago.A("class", spago.S(`bar`)),
							spago.Tag("div", 								
								spago.A("class", spago.S(`bar-item`)),
								spago.A("style", spago.S(`width: `, spago.S(c.GetProgress()), `%;`)),
								spago.T(``, spago.S(c.Current), `/(`, spago.S(c.BLE.MinID), `...`, spago.S(c.BLE.MaxID), `) -
              `, spago.S(c.GetProgress()), `%`),
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
							spago.T(`Log`),
						),
					),
					spago.Tag("div", 						
						spago.A("class", spago.S(`col-9 col-sm-12`)),
						spago.Tag("textarea", 							
							spago.A("type", spago.S(`text`)),
							spago.A("class", spago.S(`form-input input-lg`)),
							spago.A("rows", spago.S(`12`)),
							spago.A("readonly", spago.S(`true`)),
							spago.T(``, spago.S(c.Lines), ``),
						),
					),
				),
			),
		),
	)
}