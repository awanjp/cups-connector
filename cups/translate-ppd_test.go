/*
Copyright 2015 Google Inc. All rights reserved.

Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file or at
https://developers.google.com/open-source/licenses/bsd
*/
package cups

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/cups-connector/cdd"
)

func translationTest(t *testing.T, ppd string, expected *cdd.PrinterDescriptionSection) {
	description, _, _ := translatePPD(ppd)
	if !reflect.DeepEqual(expected, description) {
		e, _ := json.Marshal(expected)
		d, _ := json.Marshal(description)
		t.Logf("expected\n %s\ngot\n %s", e, d)
		t.Fail()
	}
}

func TestTrPrintingSpeed(t *testing.T) {
	ppd := `*PPD-Adobe: "4.3"
*Throughput: "30"`
	expected := &cdd.PrinterDescriptionSection{
		PrintingSpeed: &cdd.PrintingSpeed{
			[]cdd.PrintingSpeedOption{
				cdd.PrintingSpeedOption{
					SpeedPPM: 30.0,
				},
			},
		},
	}
	translationTest(t, ppd, expected)
}

func TestTrMediaSize(t *testing.T) {
	ppd := `*PPD-Adobe: "4.3"
*OpenUI *PageSize: PickOne
*DefaultPageSize: Letter
*PageSize A3/A3: ""
*PageSize ISOB5/B5 - ISO: ""
*PageSize B5/B5 - JIS: ""
*PageSize Letter/Letter: ""
*PageSize HalfLetter/5.5x8.5: ""
*CloseUI: *PageSize`
	expected := &cdd.PrinterDescriptionSection{
		MediaSize: &cdd.MediaSize{
			Option: []cdd.MediaSizeOption{
				cdd.MediaSizeOption{cdd.MediaSizeISOA3, mmToMicrons(297), mmToMicrons(420), false, false, "", "A3", cdd.NewLocalizedString("A3")},
				cdd.MediaSizeOption{cdd.MediaSizeISOB5, mmToMicrons(176), mmToMicrons(250), false, false, "", "ISOB5", cdd.NewLocalizedString("B5 (ISO)")},
				cdd.MediaSizeOption{cdd.MediaSizeJISB5, mmToMicrons(182), mmToMicrons(257), false, false, "", "B5", cdd.NewLocalizedString("B5 (JIS)")},
				cdd.MediaSizeOption{cdd.MediaSizeNALetter, inchesToMicrons(8.5), inchesToMicrons(11), false, true, "", "Letter", cdd.NewLocalizedString("Letter")},
				cdd.MediaSizeOption{cdd.MediaSizeCustom, inchesToMicrons(5.5), inchesToMicrons(8.5), false, false, "", "HalfLetter", cdd.NewLocalizedString("5.5x8.5")},
			},
		},
	}
	translationTest(t, ppd, expected)
}

func TestTrColor(t *testing.T) {
	ppd := `*PPD-Adobe: "4.3"
*OpenUI *ColorModel/Color Mode: PickOne
*DefaultColorModel: Gray
*ColorModel CMYK/Color: "(cmyk) RCsetdevicecolor"
*ColorModel Gray/Black and White: "(gray) RCsetdevicecolor"
*CloseUI: *ColorModel`
	expected := &cdd.PrinterDescriptionSection{
		Color: &cdd.Color{
			Option: []cdd.ColorOption{
				cdd.ColorOption{"ColorModelCMYK", cdd.ColorTypeStandardColor, "", false, cdd.NewLocalizedString("Color")},
				cdd.ColorOption{"ColorModelGray", cdd.ColorTypeStandardMonochrome, "", true, cdd.NewLocalizedString("Black and White")},
			},
		},
	}
	translationTest(t, ppd, expected)
}

func TestTrDuplex(t *testing.T) {
	ppd := `*PPD-Adobe: "4.3"
*OpenUI *Duplex/Duplex: PickOne
*DefaultDuplex: None
*Duplex None/Off: ""
*Duplex DuplexNoTumble/Long Edge: ""
*CloseUI: *Duplex`
	expected := &cdd.PrinterDescriptionSection{
		Duplex: &cdd.Duplex{
			Option: []cdd.DuplexOption{
				cdd.DuplexOption{cdd.DuplexNoDuplex, true},
				cdd.DuplexOption{cdd.DuplexLongEdge, false},
			},
		},
	}
	translationTest(t, ppd, expected)
}

func TestTrDPI(t *testing.T) {
	ppd := `*PPD-Adobe: "4.3"
*OpenUI *Resolution/Resolution: PickOne
*DefaultResolution: 600dpi
*Resolution 600dpi/600 dpi: ""
*Resolution 1200x600dpi/1200x600 dpi: ""
*Resolution 1200x1200dpi/1200 dpi: ""
*CloseUI: *Resolution`
	expected := &cdd.PrinterDescriptionSection{
		DPI: &cdd.DPI{
			Option: []cdd.DPIOption{
				cdd.DPIOption{600, 600, true, "", "600dpi", cdd.NewLocalizedString("600 dpi")},
				cdd.DPIOption{1200, 600, false, "", "1200x600dpi", cdd.NewLocalizedString("1200x600 dpi")},
				cdd.DPIOption{1200, 1200, false, "", "1200x1200dpi", cdd.NewLocalizedString("1200 dpi")},
			},
		},
	}
	translationTest(t, ppd, expected)
}

func TestTrInputSlot(t *testing.T) {
	ppd := `*PPD-Adobe: "4.3"
*OpenUI *OutputBin/Destination: PickOne
*OrderDependency: 210 AnySetup *OutputBin
*DefaultOutputBin: FinProof
*OutputBin Standard/Internal Tray 1: ""
*OutputBin Bin1/Internal Tray 2: ""
*OutputBin External/External Tray: ""
*CloseUI: *OutputBin`
	expected := &cdd.PrinterDescriptionSection{
		VendorCapability: &[]cdd.VendorCapability{
			cdd.VendorCapability{
				ID:                   "OutputBin",
				Type:                 cdd.VendorCapabilitySelect,
				DisplayNameLocalized: cdd.NewLocalizedString("Destination"),
				SelectCap: &cdd.SelectCapability{
					Option: []cdd.SelectCapabilityOption{
						cdd.SelectCapabilityOption{"Standard", "", true, cdd.NewLocalizedString("Internal Tray 1")},
						cdd.SelectCapabilityOption{"Bin1", "", false, cdd.NewLocalizedString("Internal Tray 2")},
						cdd.SelectCapabilityOption{"External", "", false, cdd.NewLocalizedString("External Tray")},
					},
				},
			},
		},
	}
	translationTest(t, ppd, expected)
}

func TestTrPrintQuality(t *testing.T) {
	ppd := `*PPD-Adobe: "4.3"
*OpenUI *HPPrintQuality/Print Quality: PickOne
*DefaultHPPrintQuality: FastRes1200
*HPPrintQuality FastRes1200/FastRes 1200: ""
*HPPrintQuality 600dpi/600 dpi: ""
*HPPrintQuality ProRes1200/ProRes 1200: ""
*CloseUI: *HPPrintQuality`
	expected := &cdd.PrinterDescriptionSection{
		VendorCapability: &[]cdd.VendorCapability{
			cdd.VendorCapability{
				ID:                   "HPPrintQuality",
				Type:                 cdd.VendorCapabilitySelect,
				DisplayNameLocalized: cdd.NewLocalizedString("Print Quality"),
				SelectCap: &cdd.SelectCapability{
					Option: []cdd.SelectCapabilityOption{
						cdd.SelectCapabilityOption{"FastRes1200", "", true, cdd.NewLocalizedString("FastRes 1200")},
						cdd.SelectCapabilityOption{"600dpi", "", false, cdd.NewLocalizedString("600 dpi")},
						cdd.SelectCapabilityOption{"ProRes1200", "", false, cdd.NewLocalizedString("ProRes 1200")},
					},
				},
			},
		},
	}
	translationTest(t, ppd, expected)
}

func easyModelTest(t *testing.T, input, expected string) {
	got := cleanupModel(input)
	if expected != got {
		t.Logf("expected %s got %s", expected, got)
		t.Fail()
	}
}

func TestCleanupModel(t *testing.T) {
	easyModelTest(t, "C451 PS(P)", "C451")
	easyModelTest(t, "MD-1000 Foomatic/md2k", "MD-1000")
	easyModelTest(t, "M24 Foomatic/epson (recommended)", "M24")
	easyModelTest(t, "LaserJet 2 w/PS Foomatic/Postscript (recommended)", "LaserJet 2")
	easyModelTest(t, "8445 PS2", "8445")
	easyModelTest(t, "AL-2600 PS3 v3016.103", "AL-2600")
	easyModelTest(t, "AR-163FG PS, 1.1", "AR-163FG")
	easyModelTest(t, "3212 PXL", "3212")
	easyModelTest(t, "Aficio SP C431DN PDF cups-team recommended", "Aficio SP C431DN")
	easyModelTest(t, "PIXMA Pro9000 - CUPS+Gutenprint v5.2.8-pre1", "PIXMA Pro9000")
	easyModelTest(t, "LaserJet M401dne PS A4 cups-team recommended", "LaserJet M401dne")
	easyModelTest(t, "LaserJet 4250 PS v3010.107 cups-team Letter+Duplex", "LaserJet 4250")
	easyModelTest(t, "Designjet Z5200 PostScript - PS", "Designjet Z5200")
	easyModelTest(t, "DCP-7025 BR-Script3", "DCP-7025")
	easyModelTest(t, "HL-5070DN BR-Script3J", "HL-5070DN")
	easyModelTest(t, "HL-1450 BR-Script2", "HL-1450")
	easyModelTest(t, "FS-600 (KPDL-2) Foomatic/Postscript (recommended)", "FS-600")
	easyModelTest(t, "XP-750 Series, Epson Inkjet Printer Driver (ESC/P-R) for Linux", "XP-750 Series")
	easyModelTest(t, "C5700(PS)", "C5700")
	easyModelTest(t, "OfficeJet 7400 Foomatic/hpijs (recommended) - HPLIP 0.9.7", "OfficeJet 7400")
	easyModelTest(t, "LaserJet p4015n, hpcups 3.13.9", "LaserJet p4015n")
	easyModelTest(t, "Color LaserJet 3600 hpijs, 3.13.9, requires proprietary plugin", "Color LaserJet 3600")
	easyModelTest(t, "LaserJet 4250 pcl3, hpcups 3.13.9", "LaserJet 4250")
	easyModelTest(t, "DesignJet T790 pcl, 1.0", "DesignJet T790")
}