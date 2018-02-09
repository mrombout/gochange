package changelog

// func TestRender(t *testing.T) {
// 	// arrange
// 	release020 := Release{
// 		Name: "0.2.0",
// 		Date: "2018-08-14",
// 		Added: []Entry{
// 			Entry{
// 				Description: "Some stuff.",
// 			},
// 		},
// 		PreviousRelease: &Release{
// 			Name: "NONE",
// 		},
// 	}
// 	release100 := Release{
// 		Name: "1.0.0",
// 		Date: "2018-12-28",
// 		Added: []Entry{
// 			Entry{
// 				Description: "Some stuff.",
// 			},
// 		},
// 		PreviousRelease: &release020,
// 	}

// 	currentChangelog := Changelog{
// 		URL:         "http://github.com/mrombout/gochange/",
// 		Description: "Lorum ipsum dolor sit amet consectatur.",
// 		Unreleased: Release{
// 			Name: "Unreleased",
// 			Added: []Entry{
// 				Entry{
// 					Description: "Some more stuff.",
// 				},
// 				Entry{
// 					Description: "Even more stuff.",
// 				},
// 			},
// 			Removed: []Entry{
// 				Entry{
// 					Description: "Easter egg.",
// 				},
// 				Entry{
// 					Description: "Bitcoin Miner",
// 				},
// 			},
// 			Changed: []Entry{
// 				Entry{
// 					Description: "Some things.",
// 				},
// 			},
// 		},
// 		LatestRelease: release100,
// 		Releases: []Release{
// 			release100,
// 			release020,
// 		},
// 	}

// 	// act
// 	Render(currentChangelog, os.Stdout)

// 	// assert
// 	t.Fail()
// }
