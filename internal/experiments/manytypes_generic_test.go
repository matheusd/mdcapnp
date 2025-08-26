// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

//go:build manytypesgeneric
package experiments

import "testing"
type manyAPI00000 fpFutureGeneric[string]

func (f manyAPI00000) next(s string) manyAPI00001 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00001(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00001 fpFutureGeneric[string]

func (f manyAPI00001) next(s string) manyAPI00002 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00002(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00002 fpFutureGeneric[string]

func (f manyAPI00002) next(s string) manyAPI00003 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00003(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00003 fpFutureGeneric[string]

func (f manyAPI00003) next(s string) manyAPI00004 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00004(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00004 fpFutureGeneric[string]

func (f manyAPI00004) next(s string) manyAPI00005 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00005(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00005 fpFutureGeneric[string]

func (f manyAPI00005) next(s string) manyAPI00006 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00006(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00006 fpFutureGeneric[string]

func (f manyAPI00006) next(s string) manyAPI00007 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00007(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00007 fpFutureGeneric[string]

func (f manyAPI00007) next(s string) manyAPI00008 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00008(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00008 fpFutureGeneric[string]

func (f manyAPI00008) next(s string) manyAPI00009 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00009(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00009 fpFutureGeneric[string]

func (f manyAPI00009) next(s string) manyAPI00010 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00010(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00010 fpFutureGeneric[string]

func (f manyAPI00010) next(s string) manyAPI00011 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00011(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00011 fpFutureGeneric[string]

func (f manyAPI00011) next(s string) manyAPI00012 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00012(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00012 fpFutureGeneric[string]

func (f manyAPI00012) next(s string) manyAPI00013 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00013(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00013 fpFutureGeneric[string]

func (f manyAPI00013) next(s string) manyAPI00014 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00014(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00014 fpFutureGeneric[string]

func (f manyAPI00014) next(s string) manyAPI00015 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00015(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00015 fpFutureGeneric[string]

func (f manyAPI00015) next(s string) manyAPI00016 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00016(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00016 fpFutureGeneric[string]

func (f manyAPI00016) next(s string) manyAPI00017 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00017(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00017 fpFutureGeneric[string]

func (f manyAPI00017) next(s string) manyAPI00018 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00018(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00018 fpFutureGeneric[string]

func (f manyAPI00018) next(s string) manyAPI00019 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00019(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00019 fpFutureGeneric[string]

func (f manyAPI00019) next(s string) manyAPI00020 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00020(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00020 fpFutureGeneric[string]

func (f manyAPI00020) next(s string) manyAPI00021 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00021(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00021 fpFutureGeneric[string]

func (f manyAPI00021) next(s string) manyAPI00022 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00022(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00022 fpFutureGeneric[string]

func (f manyAPI00022) next(s string) manyAPI00023 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00023(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00023 fpFutureGeneric[string]

func (f manyAPI00023) next(s string) manyAPI00024 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00024(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00024 fpFutureGeneric[string]

func (f manyAPI00024) next(s string) manyAPI00025 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00025(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00025 fpFutureGeneric[string]

func (f manyAPI00025) next(s string) manyAPI00026 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00026(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00026 fpFutureGeneric[string]

func (f manyAPI00026) next(s string) manyAPI00027 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00027(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00027 fpFutureGeneric[string]

func (f manyAPI00027) next(s string) manyAPI00028 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00028(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00028 fpFutureGeneric[string]

func (f manyAPI00028) next(s string) manyAPI00029 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00029(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00029 fpFutureGeneric[string]

func (f manyAPI00029) next(s string) manyAPI00030 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00030(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00030 fpFutureGeneric[string]

func (f manyAPI00030) next(s string) manyAPI00031 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00031(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00031 fpFutureGeneric[string]

func (f manyAPI00031) next(s string) manyAPI00032 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00032(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00032 fpFutureGeneric[string]

func (f manyAPI00032) next(s string) manyAPI00033 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00033(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00033 fpFutureGeneric[string]

func (f manyAPI00033) next(s string) manyAPI00034 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00034(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00034 fpFutureGeneric[string]

func (f manyAPI00034) next(s string) manyAPI00035 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00035(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00035 fpFutureGeneric[string]

func (f manyAPI00035) next(s string) manyAPI00036 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00036(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00036 fpFutureGeneric[string]

func (f manyAPI00036) next(s string) manyAPI00037 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00037(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00037 fpFutureGeneric[string]

func (f manyAPI00037) next(s string) manyAPI00038 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00038(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00038 fpFutureGeneric[string]

func (f manyAPI00038) next(s string) manyAPI00039 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00039(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00039 fpFutureGeneric[string]

func (f manyAPI00039) next(s string) manyAPI00040 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00040(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00040 fpFutureGeneric[string]

func (f manyAPI00040) next(s string) manyAPI00041 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00041(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00041 fpFutureGeneric[string]

func (f manyAPI00041) next(s string) manyAPI00042 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00042(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00042 fpFutureGeneric[string]

func (f manyAPI00042) next(s string) manyAPI00043 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00043(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00043 fpFutureGeneric[string]

func (f manyAPI00043) next(s string) manyAPI00044 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00044(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00044 fpFutureGeneric[string]

func (f manyAPI00044) next(s string) manyAPI00045 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00045(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00045 fpFutureGeneric[string]

func (f manyAPI00045) next(s string) manyAPI00046 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00046(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00046 fpFutureGeneric[string]

func (f manyAPI00046) next(s string) manyAPI00047 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00047(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00047 fpFutureGeneric[string]

func (f manyAPI00047) next(s string) manyAPI00048 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00048(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00048 fpFutureGeneric[string]

func (f manyAPI00048) next(s string) manyAPI00049 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00049(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00049 fpFutureGeneric[string]

func (f manyAPI00049) next(s string) manyAPI00050 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00050(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00050 fpFutureGeneric[string]

func (f manyAPI00050) next(s string) manyAPI00051 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00051(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00051 fpFutureGeneric[string]

func (f manyAPI00051) next(s string) manyAPI00052 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00052(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00052 fpFutureGeneric[string]

func (f manyAPI00052) next(s string) manyAPI00053 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00053(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00053 fpFutureGeneric[string]

func (f manyAPI00053) next(s string) manyAPI00054 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00054(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00054 fpFutureGeneric[string]

func (f manyAPI00054) next(s string) manyAPI00055 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00055(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00055 fpFutureGeneric[string]

func (f manyAPI00055) next(s string) manyAPI00056 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00056(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00056 fpFutureGeneric[string]

func (f manyAPI00056) next(s string) manyAPI00057 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00057(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00057 fpFutureGeneric[string]

func (f manyAPI00057) next(s string) manyAPI00058 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00058(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00058 fpFutureGeneric[string]

func (f manyAPI00058) next(s string) manyAPI00059 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00059(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00059 fpFutureGeneric[string]

func (f manyAPI00059) next(s string) manyAPI00060 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00060(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00060 fpFutureGeneric[string]

func (f manyAPI00060) next(s string) manyAPI00061 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00061(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00061 fpFutureGeneric[string]

func (f manyAPI00061) next(s string) manyAPI00062 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00062(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00062 fpFutureGeneric[string]

func (f manyAPI00062) next(s string) manyAPI00063 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00063(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00063 fpFutureGeneric[string]

func (f manyAPI00063) next(s string) manyAPI00064 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00064(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00064 fpFutureGeneric[string]

func (f manyAPI00064) next(s string) manyAPI00065 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00065(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00065 fpFutureGeneric[string]

func (f manyAPI00065) next(s string) manyAPI00066 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00066(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00066 fpFutureGeneric[string]

func (f manyAPI00066) next(s string) manyAPI00067 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00067(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00067 fpFutureGeneric[string]

func (f manyAPI00067) next(s string) manyAPI00068 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00068(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00068 fpFutureGeneric[string]

func (f manyAPI00068) next(s string) manyAPI00069 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00069(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00069 fpFutureGeneric[string]

func (f manyAPI00069) next(s string) manyAPI00070 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00070(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00070 fpFutureGeneric[string]

func (f manyAPI00070) next(s string) manyAPI00071 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00071(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00071 fpFutureGeneric[string]

func (f manyAPI00071) next(s string) manyAPI00072 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00072(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00072 fpFutureGeneric[string]

func (f manyAPI00072) next(s string) manyAPI00073 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00073(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00073 fpFutureGeneric[string]

func (f manyAPI00073) next(s string) manyAPI00074 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00074(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00074 fpFutureGeneric[string]

func (f manyAPI00074) next(s string) manyAPI00075 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00075(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00075 fpFutureGeneric[string]

func (f manyAPI00075) next(s string) manyAPI00076 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00076(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00076 fpFutureGeneric[string]

func (f manyAPI00076) next(s string) manyAPI00077 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00077(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00077 fpFutureGeneric[string]

func (f manyAPI00077) next(s string) manyAPI00078 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00078(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00078 fpFutureGeneric[string]

func (f manyAPI00078) next(s string) manyAPI00079 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00079(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00079 fpFutureGeneric[string]

func (f manyAPI00079) next(s string) manyAPI00080 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00080(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00080 fpFutureGeneric[string]

func (f manyAPI00080) next(s string) manyAPI00081 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00081(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00081 fpFutureGeneric[string]

func (f manyAPI00081) next(s string) manyAPI00082 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00082(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00082 fpFutureGeneric[string]

func (f manyAPI00082) next(s string) manyAPI00083 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00083(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00083 fpFutureGeneric[string]

func (f manyAPI00083) next(s string) manyAPI00084 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00084(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00084 fpFutureGeneric[string]

func (f manyAPI00084) next(s string) manyAPI00085 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00085(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00085 fpFutureGeneric[string]

func (f manyAPI00085) next(s string) manyAPI00086 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00086(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00086 fpFutureGeneric[string]

func (f manyAPI00086) next(s string) manyAPI00087 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00087(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00087 fpFutureGeneric[string]

func (f manyAPI00087) next(s string) manyAPI00088 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00088(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00088 fpFutureGeneric[string]

func (f manyAPI00088) next(s string) manyAPI00089 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00089(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00089 fpFutureGeneric[string]

func (f manyAPI00089) next(s string) manyAPI00090 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00090(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00090 fpFutureGeneric[string]

func (f manyAPI00090) next(s string) manyAPI00091 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00091(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00091 fpFutureGeneric[string]

func (f manyAPI00091) next(s string) manyAPI00092 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00092(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00092 fpFutureGeneric[string]

func (f manyAPI00092) next(s string) manyAPI00093 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00093(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00093 fpFutureGeneric[string]

func (f manyAPI00093) next(s string) manyAPI00094 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00094(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00094 fpFutureGeneric[string]

func (f manyAPI00094) next(s string) manyAPI00095 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00095(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00095 fpFutureGeneric[string]

func (f manyAPI00095) next(s string) manyAPI00096 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00096(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00096 fpFutureGeneric[string]

func (f manyAPI00096) next(s string) manyAPI00097 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00097(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00097 fpFutureGeneric[string]

func (f manyAPI00097) next(s string) manyAPI00098 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00098(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00098 fpFutureGeneric[string]

func (f manyAPI00098) next(s string) manyAPI00099 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00099(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00099 fpFutureGeneric[string]

func (f manyAPI00099) next(s string) manyAPI00100 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00100(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00100 fpFutureGeneric[string]

func (f manyAPI00100) next(s string) manyAPI00101 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00101(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00101 fpFutureGeneric[string]

func (f manyAPI00101) next(s string) manyAPI00102 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00102(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00102 fpFutureGeneric[string]

func (f manyAPI00102) next(s string) manyAPI00103 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00103(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00103 fpFutureGeneric[string]

func (f manyAPI00103) next(s string) manyAPI00104 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00104(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00104 fpFutureGeneric[string]

func (f manyAPI00104) next(s string) manyAPI00105 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00105(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00105 fpFutureGeneric[string]

func (f manyAPI00105) next(s string) manyAPI00106 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00106(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00106 fpFutureGeneric[string]

func (f manyAPI00106) next(s string) manyAPI00107 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00107(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00107 fpFutureGeneric[string]

func (f manyAPI00107) next(s string) manyAPI00108 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00108(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00108 fpFutureGeneric[string]

func (f manyAPI00108) next(s string) manyAPI00109 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00109(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00109 fpFutureGeneric[string]

func (f manyAPI00109) next(s string) manyAPI00110 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00110(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00110 fpFutureGeneric[string]

func (f manyAPI00110) next(s string) manyAPI00111 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00111(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00111 fpFutureGeneric[string]

func (f manyAPI00111) next(s string) manyAPI00112 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00112(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00112 fpFutureGeneric[string]

func (f manyAPI00112) next(s string) manyAPI00113 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00113(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00113 fpFutureGeneric[string]

func (f manyAPI00113) next(s string) manyAPI00114 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00114(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00114 fpFutureGeneric[string]

func (f manyAPI00114) next(s string) manyAPI00115 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00115(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00115 fpFutureGeneric[string]

func (f manyAPI00115) next(s string) manyAPI00116 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00116(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00116 fpFutureGeneric[string]

func (f manyAPI00116) next(s string) manyAPI00117 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00117(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00117 fpFutureGeneric[string]

func (f manyAPI00117) next(s string) manyAPI00118 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00118(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00118 fpFutureGeneric[string]

func (f manyAPI00118) next(s string) manyAPI00119 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00119(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00119 fpFutureGeneric[string]

func (f manyAPI00119) next(s string) manyAPI00120 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00120(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00120 fpFutureGeneric[string]

func (f manyAPI00120) next(s string) manyAPI00121 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00121(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00121 fpFutureGeneric[string]

func (f manyAPI00121) next(s string) manyAPI00122 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00122(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00122 fpFutureGeneric[string]

func (f manyAPI00122) next(s string) manyAPI00123 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00123(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00123 fpFutureGeneric[string]

func (f manyAPI00123) next(s string) manyAPI00124 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00124(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00124 fpFutureGeneric[string]

func (f manyAPI00124) next(s string) manyAPI00125 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00125(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00125 fpFutureGeneric[string]

func (f manyAPI00125) next(s string) manyAPI00126 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00126(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00126 fpFutureGeneric[string]

func (f manyAPI00126) next(s string) manyAPI00127 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00127(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00127 fpFutureGeneric[string]

func (f manyAPI00127) next(s string) manyAPI00128 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00128(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00128 fpFutureGeneric[string]

func (f manyAPI00128) next(s string) manyAPI00129 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00129(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00129 fpFutureGeneric[string]

func (f manyAPI00129) next(s string) manyAPI00130 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00130(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00130 fpFutureGeneric[string]

func (f manyAPI00130) next(s string) manyAPI00131 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00131(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00131 fpFutureGeneric[string]

func (f manyAPI00131) next(s string) manyAPI00132 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00132(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00132 fpFutureGeneric[string]

func (f manyAPI00132) next(s string) manyAPI00133 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00133(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00133 fpFutureGeneric[string]

func (f manyAPI00133) next(s string) manyAPI00134 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00134(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00134 fpFutureGeneric[string]

func (f manyAPI00134) next(s string) manyAPI00135 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00135(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00135 fpFutureGeneric[string]

func (f manyAPI00135) next(s string) manyAPI00136 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00136(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00136 fpFutureGeneric[string]

func (f manyAPI00136) next(s string) manyAPI00137 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00137(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00137 fpFutureGeneric[string]

func (f manyAPI00137) next(s string) manyAPI00138 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00138(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00138 fpFutureGeneric[string]

func (f manyAPI00138) next(s string) manyAPI00139 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00139(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00139 fpFutureGeneric[string]

func (f manyAPI00139) next(s string) manyAPI00140 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00140(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00140 fpFutureGeneric[string]

func (f manyAPI00140) next(s string) manyAPI00141 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00141(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00141 fpFutureGeneric[string]

func (f manyAPI00141) next(s string) manyAPI00142 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00142(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00142 fpFutureGeneric[string]

func (f manyAPI00142) next(s string) manyAPI00143 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00143(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00143 fpFutureGeneric[string]

func (f manyAPI00143) next(s string) manyAPI00144 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00144(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00144 fpFutureGeneric[string]

func (f manyAPI00144) next(s string) manyAPI00145 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00145(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00145 fpFutureGeneric[string]

func (f manyAPI00145) next(s string) manyAPI00146 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00146(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00146 fpFutureGeneric[string]

func (f manyAPI00146) next(s string) manyAPI00147 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00147(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00147 fpFutureGeneric[string]

func (f manyAPI00147) next(s string) manyAPI00148 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00148(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00148 fpFutureGeneric[string]

func (f manyAPI00148) next(s string) manyAPI00149 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00149(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00149 fpFutureGeneric[string]

func (f manyAPI00149) next(s string) manyAPI00150 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00150(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00150 fpFutureGeneric[string]

func (f manyAPI00150) next(s string) manyAPI00151 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00151(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00151 fpFutureGeneric[string]

func (f manyAPI00151) next(s string) manyAPI00152 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00152(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00152 fpFutureGeneric[string]

func (f manyAPI00152) next(s string) manyAPI00153 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00153(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00153 fpFutureGeneric[string]

func (f manyAPI00153) next(s string) manyAPI00154 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00154(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00154 fpFutureGeneric[string]

func (f manyAPI00154) next(s string) manyAPI00155 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00155(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00155 fpFutureGeneric[string]

func (f manyAPI00155) next(s string) manyAPI00156 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00156(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00156 fpFutureGeneric[string]

func (f manyAPI00156) next(s string) manyAPI00157 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00157(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00157 fpFutureGeneric[string]

func (f manyAPI00157) next(s string) manyAPI00158 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00158(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00158 fpFutureGeneric[string]

func (f manyAPI00158) next(s string) manyAPI00159 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00159(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00159 fpFutureGeneric[string]

func (f manyAPI00159) next(s string) manyAPI00160 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00160(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00160 fpFutureGeneric[string]

func (f manyAPI00160) next(s string) manyAPI00161 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00161(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00161 fpFutureGeneric[string]

func (f manyAPI00161) next(s string) manyAPI00162 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00162(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00162 fpFutureGeneric[string]

func (f manyAPI00162) next(s string) manyAPI00163 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00163(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00163 fpFutureGeneric[string]

func (f manyAPI00163) next(s string) manyAPI00164 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00164(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00164 fpFutureGeneric[string]

func (f manyAPI00164) next(s string) manyAPI00165 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00165(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00165 fpFutureGeneric[string]

func (f manyAPI00165) next(s string) manyAPI00166 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00166(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00166 fpFutureGeneric[string]

func (f manyAPI00166) next(s string) manyAPI00167 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00167(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00167 fpFutureGeneric[string]

func (f manyAPI00167) next(s string) manyAPI00168 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00168(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00168 fpFutureGeneric[string]

func (f manyAPI00168) next(s string) manyAPI00169 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00169(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00169 fpFutureGeneric[string]

func (f manyAPI00169) next(s string) manyAPI00170 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00170(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00170 fpFutureGeneric[string]

func (f manyAPI00170) next(s string) manyAPI00171 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00171(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00171 fpFutureGeneric[string]

func (f manyAPI00171) next(s string) manyAPI00172 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00172(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00172 fpFutureGeneric[string]

func (f manyAPI00172) next(s string) manyAPI00173 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00173(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00173 fpFutureGeneric[string]

func (f manyAPI00173) next(s string) manyAPI00174 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00174(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00174 fpFutureGeneric[string]

func (f manyAPI00174) next(s string) manyAPI00175 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00175(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00175 fpFutureGeneric[string]

func (f manyAPI00175) next(s string) manyAPI00176 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00176(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00176 fpFutureGeneric[string]

func (f manyAPI00176) next(s string) manyAPI00177 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00177(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00177 fpFutureGeneric[string]

func (f manyAPI00177) next(s string) manyAPI00178 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00178(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00178 fpFutureGeneric[string]

func (f manyAPI00178) next(s string) manyAPI00179 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00179(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00179 fpFutureGeneric[string]

func (f manyAPI00179) next(s string) manyAPI00180 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00180(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00180 fpFutureGeneric[string]

func (f manyAPI00180) next(s string) manyAPI00181 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00181(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00181 fpFutureGeneric[string]

func (f manyAPI00181) next(s string) manyAPI00182 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00182(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00182 fpFutureGeneric[string]

func (f manyAPI00182) next(s string) manyAPI00183 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00183(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00183 fpFutureGeneric[string]

func (f manyAPI00183) next(s string) manyAPI00184 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00184(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00184 fpFutureGeneric[string]

func (f manyAPI00184) next(s string) manyAPI00185 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00185(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00185 fpFutureGeneric[string]

func (f manyAPI00185) next(s string) manyAPI00186 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00186(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00186 fpFutureGeneric[string]

func (f manyAPI00186) next(s string) manyAPI00187 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00187(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00187 fpFutureGeneric[string]

func (f manyAPI00187) next(s string) manyAPI00188 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00188(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00188 fpFutureGeneric[string]

func (f manyAPI00188) next(s string) manyAPI00189 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00189(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00189 fpFutureGeneric[string]

func (f manyAPI00189) next(s string) manyAPI00190 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00190(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00190 fpFutureGeneric[string]

func (f manyAPI00190) next(s string) manyAPI00191 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00191(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00191 fpFutureGeneric[string]

func (f manyAPI00191) next(s string) manyAPI00192 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00192(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00192 fpFutureGeneric[string]

func (f manyAPI00192) next(s string) manyAPI00193 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00193(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00193 fpFutureGeneric[string]

func (f manyAPI00193) next(s string) manyAPI00194 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00194(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00194 fpFutureGeneric[string]

func (f manyAPI00194) next(s string) manyAPI00195 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00195(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00195 fpFutureGeneric[string]

func (f manyAPI00195) next(s string) manyAPI00196 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00196(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00196 fpFutureGeneric[string]

func (f manyAPI00196) next(s string) manyAPI00197 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00197(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00197 fpFutureGeneric[string]

func (f manyAPI00197) next(s string) manyAPI00198 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00198(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00198 fpFutureGeneric[string]

func (f manyAPI00198) next(s string) manyAPI00199 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00199(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00199 fpFutureGeneric[string]

func (f manyAPI00199) next(s string) manyAPI00200 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00200(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00200 fpFutureGeneric[string]

func (f manyAPI00200) next(s string) manyAPI00201 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00201(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00201 fpFutureGeneric[string]

func (f manyAPI00201) next(s string) manyAPI00202 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00202(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00202 fpFutureGeneric[string]

func (f manyAPI00202) next(s string) manyAPI00203 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00203(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00203 fpFutureGeneric[string]

func (f manyAPI00203) next(s string) manyAPI00204 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00204(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00204 fpFutureGeneric[string]

func (f manyAPI00204) next(s string) manyAPI00205 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00205(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00205 fpFutureGeneric[string]

func (f manyAPI00205) next(s string) manyAPI00206 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00206(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00206 fpFutureGeneric[string]

func (f manyAPI00206) next(s string) manyAPI00207 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00207(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00207 fpFutureGeneric[string]

func (f manyAPI00207) next(s string) manyAPI00208 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00208(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00208 fpFutureGeneric[string]

func (f manyAPI00208) next(s string) manyAPI00209 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00209(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00209 fpFutureGeneric[string]

func (f manyAPI00209) next(s string) manyAPI00210 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00210(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00210 fpFutureGeneric[string]

func (f manyAPI00210) next(s string) manyAPI00211 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00211(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00211 fpFutureGeneric[string]

func (f manyAPI00211) next(s string) manyAPI00212 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00212(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00212 fpFutureGeneric[string]

func (f manyAPI00212) next(s string) manyAPI00213 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00213(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00213 fpFutureGeneric[string]

func (f manyAPI00213) next(s string) manyAPI00214 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00214(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00214 fpFutureGeneric[string]

func (f manyAPI00214) next(s string) manyAPI00215 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00215(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00215 fpFutureGeneric[string]

func (f manyAPI00215) next(s string) manyAPI00216 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00216(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00216 fpFutureGeneric[string]

func (f manyAPI00216) next(s string) manyAPI00217 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00217(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00217 fpFutureGeneric[string]

func (f manyAPI00217) next(s string) manyAPI00218 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00218(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00218 fpFutureGeneric[string]

func (f manyAPI00218) next(s string) manyAPI00219 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00219(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00219 fpFutureGeneric[string]

func (f manyAPI00219) next(s string) manyAPI00220 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00220(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00220 fpFutureGeneric[string]

func (f manyAPI00220) next(s string) manyAPI00221 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00221(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00221 fpFutureGeneric[string]

func (f manyAPI00221) next(s string) manyAPI00222 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00222(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00222 fpFutureGeneric[string]

func (f manyAPI00222) next(s string) manyAPI00223 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00223(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00223 fpFutureGeneric[string]

func (f manyAPI00223) next(s string) manyAPI00224 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00224(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00224 fpFutureGeneric[string]

func (f manyAPI00224) next(s string) manyAPI00225 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00225(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00225 fpFutureGeneric[string]

func (f manyAPI00225) next(s string) manyAPI00226 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00226(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00226 fpFutureGeneric[string]

func (f manyAPI00226) next(s string) manyAPI00227 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00227(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00227 fpFutureGeneric[string]

func (f manyAPI00227) next(s string) manyAPI00228 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00228(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00228 fpFutureGeneric[string]

func (f manyAPI00228) next(s string) manyAPI00229 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00229(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00229 fpFutureGeneric[string]

func (f manyAPI00229) next(s string) manyAPI00230 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00230(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00230 fpFutureGeneric[string]

func (f manyAPI00230) next(s string) manyAPI00231 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00231(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00231 fpFutureGeneric[string]

func (f manyAPI00231) next(s string) manyAPI00232 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00232(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00232 fpFutureGeneric[string]

func (f manyAPI00232) next(s string) manyAPI00233 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00233(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00233 fpFutureGeneric[string]

func (f manyAPI00233) next(s string) manyAPI00234 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00234(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00234 fpFutureGeneric[string]

func (f manyAPI00234) next(s string) manyAPI00235 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00235(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00235 fpFutureGeneric[string]

func (f manyAPI00235) next(s string) manyAPI00236 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00236(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00236 fpFutureGeneric[string]

func (f manyAPI00236) next(s string) manyAPI00237 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00237(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00237 fpFutureGeneric[string]

func (f manyAPI00237) next(s string) manyAPI00238 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00238(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00238 fpFutureGeneric[string]

func (f manyAPI00238) next(s string) manyAPI00239 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00239(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00239 fpFutureGeneric[string]

func (f manyAPI00239) next(s string) manyAPI00240 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00240(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00240 fpFutureGeneric[string]

func (f manyAPI00240) next(s string) manyAPI00241 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00241(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00241 fpFutureGeneric[string]

func (f manyAPI00241) next(s string) manyAPI00242 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00242(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00242 fpFutureGeneric[string]

func (f manyAPI00242) next(s string) manyAPI00243 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00243(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00243 fpFutureGeneric[string]

func (f manyAPI00243) next(s string) manyAPI00244 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00244(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00244 fpFutureGeneric[string]

func (f manyAPI00244) next(s string) manyAPI00245 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00245(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00245 fpFutureGeneric[string]

func (f manyAPI00245) next(s string) manyAPI00246 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00246(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00246 fpFutureGeneric[string]

func (f manyAPI00246) next(s string) manyAPI00247 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00247(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00247 fpFutureGeneric[string]

func (f manyAPI00247) next(s string) manyAPI00248 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00248(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00248 fpFutureGeneric[string]

func (f manyAPI00248) next(s string) manyAPI00249 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00249(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00249 fpFutureGeneric[string]

func (f manyAPI00249) next(s string) manyAPI00250 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00250(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00250 fpFutureGeneric[string]

func (f manyAPI00250) next(s string) manyAPI00251 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00251(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00251 fpFutureGeneric[string]

func (f manyAPI00251) next(s string) manyAPI00252 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00252(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00252 fpFutureGeneric[string]

func (f manyAPI00252) next(s string) manyAPI00253 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00253(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00253 fpFutureGeneric[string]

func (f manyAPI00253) next(s string) manyAPI00254 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00254(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00254 fpFutureGeneric[string]

func (f manyAPI00254) next(s string) manyAPI00255 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00255(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00255 fpFutureGeneric[string]

func (f manyAPI00255) next(s string) manyAPI00256 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00256(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00256 fpFutureGeneric[string]

func (f manyAPI00256) next(s string) manyAPI00257 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00257(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00257 fpFutureGeneric[string]

func (f manyAPI00257) next(s string) manyAPI00258 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00258(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00258 fpFutureGeneric[string]

func (f manyAPI00258) next(s string) manyAPI00259 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00259(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00259 fpFutureGeneric[string]

func (f manyAPI00259) next(s string) manyAPI00260 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00260(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00260 fpFutureGeneric[string]

func (f manyAPI00260) next(s string) manyAPI00261 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00261(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00261 fpFutureGeneric[string]

func (f manyAPI00261) next(s string) manyAPI00262 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00262(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00262 fpFutureGeneric[string]

func (f manyAPI00262) next(s string) manyAPI00263 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00263(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00263 fpFutureGeneric[string]

func (f manyAPI00263) next(s string) manyAPI00264 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00264(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00264 fpFutureGeneric[string]

func (f manyAPI00264) next(s string) manyAPI00265 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00265(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00265 fpFutureGeneric[string]

func (f manyAPI00265) next(s string) manyAPI00266 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00266(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00266 fpFutureGeneric[string]

func (f manyAPI00266) next(s string) manyAPI00267 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00267(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00267 fpFutureGeneric[string]

func (f manyAPI00267) next(s string) manyAPI00268 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00268(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00268 fpFutureGeneric[string]

func (f manyAPI00268) next(s string) manyAPI00269 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00269(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00269 fpFutureGeneric[string]

func (f manyAPI00269) next(s string) manyAPI00270 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00270(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00270 fpFutureGeneric[string]

func (f manyAPI00270) next(s string) manyAPI00271 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00271(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00271 fpFutureGeneric[string]

func (f manyAPI00271) next(s string) manyAPI00272 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00272(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00272 fpFutureGeneric[string]

func (f manyAPI00272) next(s string) manyAPI00273 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00273(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00273 fpFutureGeneric[string]

func (f manyAPI00273) next(s string) manyAPI00274 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00274(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00274 fpFutureGeneric[string]

func (f manyAPI00274) next(s string) manyAPI00275 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00275(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00275 fpFutureGeneric[string]

func (f manyAPI00275) next(s string) manyAPI00276 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00276(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00276 fpFutureGeneric[string]

func (f manyAPI00276) next(s string) manyAPI00277 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00277(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00277 fpFutureGeneric[string]

func (f manyAPI00277) next(s string) manyAPI00278 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00278(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00278 fpFutureGeneric[string]

func (f manyAPI00278) next(s string) manyAPI00279 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00279(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00279 fpFutureGeneric[string]

func (f manyAPI00279) next(s string) manyAPI00280 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00280(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00280 fpFutureGeneric[string]

func (f manyAPI00280) next(s string) manyAPI00281 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00281(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00281 fpFutureGeneric[string]

func (f manyAPI00281) next(s string) manyAPI00282 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00282(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00282 fpFutureGeneric[string]

func (f manyAPI00282) next(s string) manyAPI00283 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00283(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00283 fpFutureGeneric[string]

func (f manyAPI00283) next(s string) manyAPI00284 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00284(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00284 fpFutureGeneric[string]

func (f manyAPI00284) next(s string) manyAPI00285 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00285(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00285 fpFutureGeneric[string]

func (f manyAPI00285) next(s string) manyAPI00286 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00286(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00286 fpFutureGeneric[string]

func (f manyAPI00286) next(s string) manyAPI00287 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00287(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00287 fpFutureGeneric[string]

func (f manyAPI00287) next(s string) manyAPI00288 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00288(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00288 fpFutureGeneric[string]

func (f manyAPI00288) next(s string) manyAPI00289 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00289(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00289 fpFutureGeneric[string]

func (f manyAPI00289) next(s string) manyAPI00290 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00290(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00290 fpFutureGeneric[string]

func (f manyAPI00290) next(s string) manyAPI00291 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00291(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00291 fpFutureGeneric[string]

func (f manyAPI00291) next(s string) manyAPI00292 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00292(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00292 fpFutureGeneric[string]

func (f manyAPI00292) next(s string) manyAPI00293 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00293(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00293 fpFutureGeneric[string]

func (f manyAPI00293) next(s string) manyAPI00294 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00294(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00294 fpFutureGeneric[string]

func (f manyAPI00294) next(s string) manyAPI00295 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00295(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00295 fpFutureGeneric[string]

func (f manyAPI00295) next(s string) manyAPI00296 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00296(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00296 fpFutureGeneric[string]

func (f manyAPI00296) next(s string) manyAPI00297 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00297(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00297 fpFutureGeneric[string]

func (f manyAPI00297) next(s string) manyAPI00298 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00298(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00298 fpFutureGeneric[string]

func (f manyAPI00298) next(s string) manyAPI00299 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00299(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00299 fpFutureGeneric[string]

func (f manyAPI00299) next(s string) manyAPI00300 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00300(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00300 fpFutureGeneric[string]

func (f manyAPI00300) next(s string) manyAPI00301 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00301(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00301 fpFutureGeneric[string]

func (f manyAPI00301) next(s string) manyAPI00302 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00302(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00302 fpFutureGeneric[string]

func (f manyAPI00302) next(s string) manyAPI00303 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00303(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00303 fpFutureGeneric[string]

func (f manyAPI00303) next(s string) manyAPI00304 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00304(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00304 fpFutureGeneric[string]

func (f manyAPI00304) next(s string) manyAPI00305 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00305(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00305 fpFutureGeneric[string]

func (f manyAPI00305) next(s string) manyAPI00306 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00306(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00306 fpFutureGeneric[string]

func (f manyAPI00306) next(s string) manyAPI00307 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00307(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00307 fpFutureGeneric[string]

func (f manyAPI00307) next(s string) manyAPI00308 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00308(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00308 fpFutureGeneric[string]

func (f manyAPI00308) next(s string) manyAPI00309 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00309(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00309 fpFutureGeneric[string]

func (f manyAPI00309) next(s string) manyAPI00310 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00310(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00310 fpFutureGeneric[string]

func (f manyAPI00310) next(s string) manyAPI00311 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00311(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00311 fpFutureGeneric[string]

func (f manyAPI00311) next(s string) manyAPI00312 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00312(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00312 fpFutureGeneric[string]

func (f manyAPI00312) next(s string) manyAPI00313 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00313(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00313 fpFutureGeneric[string]

func (f manyAPI00313) next(s string) manyAPI00314 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00314(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00314 fpFutureGeneric[string]

func (f manyAPI00314) next(s string) manyAPI00315 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00315(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00315 fpFutureGeneric[string]

func (f manyAPI00315) next(s string) manyAPI00316 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00316(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00316 fpFutureGeneric[string]

func (f manyAPI00316) next(s string) manyAPI00317 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00317(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00317 fpFutureGeneric[string]

func (f manyAPI00317) next(s string) manyAPI00318 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00318(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00318 fpFutureGeneric[string]

func (f manyAPI00318) next(s string) manyAPI00319 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00319(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00319 fpFutureGeneric[string]

func (f manyAPI00319) next(s string) manyAPI00320 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00320(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00320 fpFutureGeneric[string]

func (f manyAPI00320) next(s string) manyAPI00321 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00321(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00321 fpFutureGeneric[string]

func (f manyAPI00321) next(s string) manyAPI00322 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00322(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00322 fpFutureGeneric[string]

func (f manyAPI00322) next(s string) manyAPI00323 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00323(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00323 fpFutureGeneric[string]

func (f manyAPI00323) next(s string) manyAPI00324 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00324(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00324 fpFutureGeneric[string]

func (f manyAPI00324) next(s string) manyAPI00325 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00325(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00325 fpFutureGeneric[string]

func (f manyAPI00325) next(s string) manyAPI00326 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00326(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00326 fpFutureGeneric[string]

func (f manyAPI00326) next(s string) manyAPI00327 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00327(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00327 fpFutureGeneric[string]

func (f manyAPI00327) next(s string) manyAPI00328 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00328(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00328 fpFutureGeneric[string]

func (f manyAPI00328) next(s string) manyAPI00329 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00329(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00329 fpFutureGeneric[string]

func (f manyAPI00329) next(s string) manyAPI00330 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00330(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00330 fpFutureGeneric[string]

func (f manyAPI00330) next(s string) manyAPI00331 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00331(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00331 fpFutureGeneric[string]

func (f manyAPI00331) next(s string) manyAPI00332 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00332(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00332 fpFutureGeneric[string]

func (f manyAPI00332) next(s string) manyAPI00333 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00333(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00333 fpFutureGeneric[string]

func (f manyAPI00333) next(s string) manyAPI00334 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00334(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00334 fpFutureGeneric[string]

func (f manyAPI00334) next(s string) manyAPI00335 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00335(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00335 fpFutureGeneric[string]

func (f manyAPI00335) next(s string) manyAPI00336 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00336(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00336 fpFutureGeneric[string]

func (f manyAPI00336) next(s string) manyAPI00337 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00337(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00337 fpFutureGeneric[string]

func (f manyAPI00337) next(s string) manyAPI00338 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00338(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00338 fpFutureGeneric[string]

func (f manyAPI00338) next(s string) manyAPI00339 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00339(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00339 fpFutureGeneric[string]

func (f manyAPI00339) next(s string) manyAPI00340 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00340(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00340 fpFutureGeneric[string]

func (f manyAPI00340) next(s string) manyAPI00341 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00341(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00341 fpFutureGeneric[string]

func (f manyAPI00341) next(s string) manyAPI00342 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00342(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00342 fpFutureGeneric[string]

func (f manyAPI00342) next(s string) manyAPI00343 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00343(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00343 fpFutureGeneric[string]

func (f manyAPI00343) next(s string) manyAPI00344 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00344(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00344 fpFutureGeneric[string]

func (f manyAPI00344) next(s string) manyAPI00345 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00345(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00345 fpFutureGeneric[string]

func (f manyAPI00345) next(s string) manyAPI00346 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00346(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00346 fpFutureGeneric[string]

func (f manyAPI00346) next(s string) manyAPI00347 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00347(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00347 fpFutureGeneric[string]

func (f manyAPI00347) next(s string) manyAPI00348 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00348(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00348 fpFutureGeneric[string]

func (f manyAPI00348) next(s string) manyAPI00349 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00349(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00349 fpFutureGeneric[string]

func (f manyAPI00349) next(s string) manyAPI00350 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00350(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00350 fpFutureGeneric[string]

func (f manyAPI00350) next(s string) manyAPI00351 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00351(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00351 fpFutureGeneric[string]

func (f manyAPI00351) next(s string) manyAPI00352 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00352(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00352 fpFutureGeneric[string]

func (f manyAPI00352) next(s string) manyAPI00353 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00353(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00353 fpFutureGeneric[string]

func (f manyAPI00353) next(s string) manyAPI00354 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00354(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00354 fpFutureGeneric[string]

func (f manyAPI00354) next(s string) manyAPI00355 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00355(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00355 fpFutureGeneric[string]

func (f manyAPI00355) next(s string) manyAPI00356 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00356(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00356 fpFutureGeneric[string]

func (f manyAPI00356) next(s string) manyAPI00357 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00357(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00357 fpFutureGeneric[string]

func (f manyAPI00357) next(s string) manyAPI00358 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00358(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00358 fpFutureGeneric[string]

func (f manyAPI00358) next(s string) manyAPI00359 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00359(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00359 fpFutureGeneric[string]

func (f manyAPI00359) next(s string) manyAPI00360 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00360(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00360 fpFutureGeneric[string]

func (f manyAPI00360) next(s string) manyAPI00361 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00361(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00361 fpFutureGeneric[string]

func (f manyAPI00361) next(s string) manyAPI00362 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00362(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00362 fpFutureGeneric[string]

func (f manyAPI00362) next(s string) manyAPI00363 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00363(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00363 fpFutureGeneric[string]

func (f manyAPI00363) next(s string) manyAPI00364 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00364(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00364 fpFutureGeneric[string]

func (f manyAPI00364) next(s string) manyAPI00365 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00365(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00365 fpFutureGeneric[string]

func (f manyAPI00365) next(s string) manyAPI00366 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00366(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00366 fpFutureGeneric[string]

func (f manyAPI00366) next(s string) manyAPI00367 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00367(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00367 fpFutureGeneric[string]

func (f manyAPI00367) next(s string) manyAPI00368 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00368(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00368 fpFutureGeneric[string]

func (f manyAPI00368) next(s string) manyAPI00369 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00369(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00369 fpFutureGeneric[string]

func (f manyAPI00369) next(s string) manyAPI00370 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00370(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00370 fpFutureGeneric[string]

func (f manyAPI00370) next(s string) manyAPI00371 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00371(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00371 fpFutureGeneric[string]

func (f manyAPI00371) next(s string) manyAPI00372 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00372(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00372 fpFutureGeneric[string]

func (f manyAPI00372) next(s string) manyAPI00373 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00373(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00373 fpFutureGeneric[string]

func (f manyAPI00373) next(s string) manyAPI00374 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00374(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00374 fpFutureGeneric[string]

func (f manyAPI00374) next(s string) manyAPI00375 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00375(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00375 fpFutureGeneric[string]

func (f manyAPI00375) next(s string) manyAPI00376 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00376(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00376 fpFutureGeneric[string]

func (f manyAPI00376) next(s string) manyAPI00377 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00377(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00377 fpFutureGeneric[string]

func (f manyAPI00377) next(s string) manyAPI00378 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00378(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00378 fpFutureGeneric[string]

func (f manyAPI00378) next(s string) manyAPI00379 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00379(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00379 fpFutureGeneric[string]

func (f manyAPI00379) next(s string) manyAPI00380 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00380(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00380 fpFutureGeneric[string]

func (f manyAPI00380) next(s string) manyAPI00381 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00381(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00381 fpFutureGeneric[string]

func (f manyAPI00381) next(s string) manyAPI00382 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00382(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00382 fpFutureGeneric[string]

func (f manyAPI00382) next(s string) manyAPI00383 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00383(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00383 fpFutureGeneric[string]

func (f manyAPI00383) next(s string) manyAPI00384 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00384(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00384 fpFutureGeneric[string]

func (f manyAPI00384) next(s string) manyAPI00385 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00385(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00385 fpFutureGeneric[string]

func (f manyAPI00385) next(s string) manyAPI00386 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00386(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00386 fpFutureGeneric[string]

func (f manyAPI00386) next(s string) manyAPI00387 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00387(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00387 fpFutureGeneric[string]

func (f manyAPI00387) next(s string) manyAPI00388 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00388(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00388 fpFutureGeneric[string]

func (f manyAPI00388) next(s string) manyAPI00389 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00389(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00389 fpFutureGeneric[string]

func (f manyAPI00389) next(s string) manyAPI00390 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00390(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00390 fpFutureGeneric[string]

func (f manyAPI00390) next(s string) manyAPI00391 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00391(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00391 fpFutureGeneric[string]

func (f manyAPI00391) next(s string) manyAPI00392 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00392(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00392 fpFutureGeneric[string]

func (f manyAPI00392) next(s string) manyAPI00393 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00393(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00393 fpFutureGeneric[string]

func (f manyAPI00393) next(s string) manyAPI00394 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00394(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00394 fpFutureGeneric[string]

func (f manyAPI00394) next(s string) manyAPI00395 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00395(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00395 fpFutureGeneric[string]

func (f manyAPI00395) next(s string) manyAPI00396 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00396(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00396 fpFutureGeneric[string]

func (f manyAPI00396) next(s string) manyAPI00397 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00397(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00397 fpFutureGeneric[string]

func (f manyAPI00397) next(s string) manyAPI00398 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00398(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00398 fpFutureGeneric[string]

func (f manyAPI00398) next(s string) manyAPI00399 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00399(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00399 fpFutureGeneric[string]

func (f manyAPI00399) next(s string) manyAPI00400 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00400(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00400 fpFutureGeneric[string]

func (f manyAPI00400) next(s string) manyAPI00401 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00401(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00401 fpFutureGeneric[string]

func (f manyAPI00401) next(s string) manyAPI00402 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00402(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00402 fpFutureGeneric[string]

func (f manyAPI00402) next(s string) manyAPI00403 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00403(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00403 fpFutureGeneric[string]

func (f manyAPI00403) next(s string) manyAPI00404 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00404(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00404 fpFutureGeneric[string]

func (f manyAPI00404) next(s string) manyAPI00405 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00405(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00405 fpFutureGeneric[string]

func (f manyAPI00405) next(s string) manyAPI00406 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00406(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00406 fpFutureGeneric[string]

func (f manyAPI00406) next(s string) manyAPI00407 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00407(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00407 fpFutureGeneric[string]

func (f manyAPI00407) next(s string) manyAPI00408 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00408(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00408 fpFutureGeneric[string]

func (f manyAPI00408) next(s string) manyAPI00409 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00409(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00409 fpFutureGeneric[string]

func (f manyAPI00409) next(s string) manyAPI00410 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00410(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00410 fpFutureGeneric[string]

func (f manyAPI00410) next(s string) manyAPI00411 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00411(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00411 fpFutureGeneric[string]

func (f manyAPI00411) next(s string) manyAPI00412 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00412(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00412 fpFutureGeneric[string]

func (f manyAPI00412) next(s string) manyAPI00413 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00413(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00413 fpFutureGeneric[string]

func (f manyAPI00413) next(s string) manyAPI00414 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00414(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00414 fpFutureGeneric[string]

func (f manyAPI00414) next(s string) manyAPI00415 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00415(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00415 fpFutureGeneric[string]

func (f manyAPI00415) next(s string) manyAPI00416 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00416(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00416 fpFutureGeneric[string]

func (f manyAPI00416) next(s string) manyAPI00417 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00417(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00417 fpFutureGeneric[string]

func (f manyAPI00417) next(s string) manyAPI00418 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00418(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00418 fpFutureGeneric[string]

func (f manyAPI00418) next(s string) manyAPI00419 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00419(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00419 fpFutureGeneric[string]

func (f manyAPI00419) next(s string) manyAPI00420 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00420(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00420 fpFutureGeneric[string]

func (f manyAPI00420) next(s string) manyAPI00421 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00421(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00421 fpFutureGeneric[string]

func (f manyAPI00421) next(s string) manyAPI00422 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00422(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00422 fpFutureGeneric[string]

func (f manyAPI00422) next(s string) manyAPI00423 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00423(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00423 fpFutureGeneric[string]

func (f manyAPI00423) next(s string) manyAPI00424 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00424(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00424 fpFutureGeneric[string]

func (f manyAPI00424) next(s string) manyAPI00425 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00425(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00425 fpFutureGeneric[string]

func (f manyAPI00425) next(s string) manyAPI00426 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00426(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00426 fpFutureGeneric[string]

func (f manyAPI00426) next(s string) manyAPI00427 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00427(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00427 fpFutureGeneric[string]

func (f manyAPI00427) next(s string) manyAPI00428 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00428(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00428 fpFutureGeneric[string]

func (f manyAPI00428) next(s string) manyAPI00429 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00429(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00429 fpFutureGeneric[string]

func (f manyAPI00429) next(s string) manyAPI00430 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00430(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00430 fpFutureGeneric[string]

func (f manyAPI00430) next(s string) manyAPI00431 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00431(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00431 fpFutureGeneric[string]

func (f manyAPI00431) next(s string) manyAPI00432 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00432(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00432 fpFutureGeneric[string]

func (f manyAPI00432) next(s string) manyAPI00433 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00433(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00433 fpFutureGeneric[string]

func (f manyAPI00433) next(s string) manyAPI00434 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00434(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00434 fpFutureGeneric[string]

func (f manyAPI00434) next(s string) manyAPI00435 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00435(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00435 fpFutureGeneric[string]

func (f manyAPI00435) next(s string) manyAPI00436 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00436(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00436 fpFutureGeneric[string]

func (f manyAPI00436) next(s string) manyAPI00437 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00437(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00437 fpFutureGeneric[string]

func (f manyAPI00437) next(s string) manyAPI00438 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00438(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00438 fpFutureGeneric[string]

func (f manyAPI00438) next(s string) manyAPI00439 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00439(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00439 fpFutureGeneric[string]

func (f manyAPI00439) next(s string) manyAPI00440 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00440(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00440 fpFutureGeneric[string]

func (f manyAPI00440) next(s string) manyAPI00441 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00441(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00441 fpFutureGeneric[string]

func (f manyAPI00441) next(s string) manyAPI00442 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00442(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00442 fpFutureGeneric[string]

func (f manyAPI00442) next(s string) manyAPI00443 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00443(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00443 fpFutureGeneric[string]

func (f manyAPI00443) next(s string) manyAPI00444 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00444(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00444 fpFutureGeneric[string]

func (f manyAPI00444) next(s string) manyAPI00445 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00445(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00445 fpFutureGeneric[string]

func (f manyAPI00445) next(s string) manyAPI00446 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00446(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00446 fpFutureGeneric[string]

func (f manyAPI00446) next(s string) manyAPI00447 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00447(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00447 fpFutureGeneric[string]

func (f manyAPI00447) next(s string) manyAPI00448 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00448(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00448 fpFutureGeneric[string]

func (f manyAPI00448) next(s string) manyAPI00449 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00449(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00449 fpFutureGeneric[string]

func (f manyAPI00449) next(s string) manyAPI00450 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00450(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00450 fpFutureGeneric[string]

func (f manyAPI00450) next(s string) manyAPI00451 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00451(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00451 fpFutureGeneric[string]

func (f manyAPI00451) next(s string) manyAPI00452 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00452(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00452 fpFutureGeneric[string]

func (f manyAPI00452) next(s string) manyAPI00453 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00453(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00453 fpFutureGeneric[string]

func (f manyAPI00453) next(s string) manyAPI00454 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00454(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00454 fpFutureGeneric[string]

func (f manyAPI00454) next(s string) manyAPI00455 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00455(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00455 fpFutureGeneric[string]

func (f manyAPI00455) next(s string) manyAPI00456 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00456(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00456 fpFutureGeneric[string]

func (f manyAPI00456) next(s string) manyAPI00457 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00457(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00457 fpFutureGeneric[string]

func (f manyAPI00457) next(s string) manyAPI00458 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00458(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00458 fpFutureGeneric[string]

func (f manyAPI00458) next(s string) manyAPI00459 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00459(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00459 fpFutureGeneric[string]

func (f manyAPI00459) next(s string) manyAPI00460 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00460(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00460 fpFutureGeneric[string]

func (f manyAPI00460) next(s string) manyAPI00461 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00461(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00461 fpFutureGeneric[string]

func (f manyAPI00461) next(s string) manyAPI00462 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00462(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00462 fpFutureGeneric[string]

func (f manyAPI00462) next(s string) manyAPI00463 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00463(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00463 fpFutureGeneric[string]

func (f manyAPI00463) next(s string) manyAPI00464 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00464(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00464 fpFutureGeneric[string]

func (f manyAPI00464) next(s string) manyAPI00465 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00465(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00465 fpFutureGeneric[string]

func (f manyAPI00465) next(s string) manyAPI00466 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00466(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00466 fpFutureGeneric[string]

func (f manyAPI00466) next(s string) manyAPI00467 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00467(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00467 fpFutureGeneric[string]

func (f manyAPI00467) next(s string) manyAPI00468 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00468(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00468 fpFutureGeneric[string]

func (f manyAPI00468) next(s string) manyAPI00469 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00469(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00469 fpFutureGeneric[string]

func (f manyAPI00469) next(s string) manyAPI00470 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00470(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00470 fpFutureGeneric[string]

func (f manyAPI00470) next(s string) manyAPI00471 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00471(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00471 fpFutureGeneric[string]

func (f manyAPI00471) next(s string) manyAPI00472 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00472(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00472 fpFutureGeneric[string]

func (f manyAPI00472) next(s string) manyAPI00473 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00473(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00473 fpFutureGeneric[string]

func (f manyAPI00473) next(s string) manyAPI00474 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00474(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00474 fpFutureGeneric[string]

func (f manyAPI00474) next(s string) manyAPI00475 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00475(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00475 fpFutureGeneric[string]

func (f manyAPI00475) next(s string) manyAPI00476 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00476(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00476 fpFutureGeneric[string]

func (f manyAPI00476) next(s string) manyAPI00477 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00477(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00477 fpFutureGeneric[string]

func (f manyAPI00477) next(s string) manyAPI00478 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00478(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00478 fpFutureGeneric[string]

func (f manyAPI00478) next(s string) manyAPI00479 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00479(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00479 fpFutureGeneric[string]

func (f manyAPI00479) next(s string) manyAPI00480 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00480(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00480 fpFutureGeneric[string]

func (f manyAPI00480) next(s string) manyAPI00481 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00481(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00481 fpFutureGeneric[string]

func (f manyAPI00481) next(s string) manyAPI00482 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00482(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00482 fpFutureGeneric[string]

func (f manyAPI00482) next(s string) manyAPI00483 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00483(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00483 fpFutureGeneric[string]

func (f manyAPI00483) next(s string) manyAPI00484 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00484(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00484 fpFutureGeneric[string]

func (f manyAPI00484) next(s string) manyAPI00485 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00485(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00485 fpFutureGeneric[string]

func (f manyAPI00485) next(s string) manyAPI00486 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00486(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00486 fpFutureGeneric[string]

func (f manyAPI00486) next(s string) manyAPI00487 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00487(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00487 fpFutureGeneric[string]

func (f manyAPI00487) next(s string) manyAPI00488 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00488(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00488 fpFutureGeneric[string]

func (f manyAPI00488) next(s string) manyAPI00489 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00489(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00489 fpFutureGeneric[string]

func (f manyAPI00489) next(s string) manyAPI00490 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00490(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00490 fpFutureGeneric[string]

func (f manyAPI00490) next(s string) manyAPI00491 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00491(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00491 fpFutureGeneric[string]

func (f manyAPI00491) next(s string) manyAPI00492 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00492(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00492 fpFutureGeneric[string]

func (f manyAPI00492) next(s string) manyAPI00493 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00493(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00493 fpFutureGeneric[string]

func (f manyAPI00493) next(s string) manyAPI00494 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00494(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00494 fpFutureGeneric[string]

func (f manyAPI00494) next(s string) manyAPI00495 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00495(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00495 fpFutureGeneric[string]

func (f manyAPI00495) next(s string) manyAPI00496 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00496(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00496 fpFutureGeneric[string]

func (f manyAPI00496) next(s string) manyAPI00497 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00497(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00497 fpFutureGeneric[string]

func (f manyAPI00497) next(s string) manyAPI00498 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00498(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00498 fpFutureGeneric[string]

func (f manyAPI00498) next(s string) manyAPI00499 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00499(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00499 fpFutureGeneric[string]

func (f manyAPI00499) next(s string) manyAPI00500 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00500(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00500 fpFutureGeneric[string]

func (f manyAPI00500) next(s string) manyAPI00501 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00501(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00501 fpFutureGeneric[string]

func (f manyAPI00501) next(s string) manyAPI00502 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00502(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00502 fpFutureGeneric[string]

func (f manyAPI00502) next(s string) manyAPI00503 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00503(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00503 fpFutureGeneric[string]

func (f manyAPI00503) next(s string) manyAPI00504 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00504(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00504 fpFutureGeneric[string]

func (f manyAPI00504) next(s string) manyAPI00505 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00505(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00505 fpFutureGeneric[string]

func (f manyAPI00505) next(s string) manyAPI00506 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00506(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00506 fpFutureGeneric[string]

func (f manyAPI00506) next(s string) manyAPI00507 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00507(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00507 fpFutureGeneric[string]

func (f manyAPI00507) next(s string) manyAPI00508 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00508(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00508 fpFutureGeneric[string]

func (f manyAPI00508) next(s string) manyAPI00509 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00509(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00509 fpFutureGeneric[string]

func (f manyAPI00509) next(s string) manyAPI00510 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00510(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00510 fpFutureGeneric[string]

func (f manyAPI00510) next(s string) manyAPI00511 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00511(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00511 fpFutureGeneric[string]

func (f manyAPI00511) next(s string) manyAPI00512 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00512(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00512 fpFutureGeneric[string]

func (f manyAPI00512) next(s string) manyAPI00513 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00513(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00513 fpFutureGeneric[string]

func (f manyAPI00513) next(s string) manyAPI00514 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00514(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00514 fpFutureGeneric[string]

func (f manyAPI00514) next(s string) manyAPI00515 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00515(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00515 fpFutureGeneric[string]

func (f manyAPI00515) next(s string) manyAPI00516 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00516(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00516 fpFutureGeneric[string]

func (f manyAPI00516) next(s string) manyAPI00517 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00517(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00517 fpFutureGeneric[string]

func (f manyAPI00517) next(s string) manyAPI00518 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00518(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00518 fpFutureGeneric[string]

func (f manyAPI00518) next(s string) manyAPI00519 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00519(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00519 fpFutureGeneric[string]

func (f manyAPI00519) next(s string) manyAPI00520 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00520(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00520 fpFutureGeneric[string]

func (f manyAPI00520) next(s string) manyAPI00521 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00521(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00521 fpFutureGeneric[string]

func (f manyAPI00521) next(s string) manyAPI00522 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00522(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00522 fpFutureGeneric[string]

func (f manyAPI00522) next(s string) manyAPI00523 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00523(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00523 fpFutureGeneric[string]

func (f manyAPI00523) next(s string) manyAPI00524 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00524(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00524 fpFutureGeneric[string]

func (f manyAPI00524) next(s string) manyAPI00525 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00525(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00525 fpFutureGeneric[string]

func (f manyAPI00525) next(s string) manyAPI00526 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00526(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00526 fpFutureGeneric[string]

func (f manyAPI00526) next(s string) manyAPI00527 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00527(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00527 fpFutureGeneric[string]

func (f manyAPI00527) next(s string) manyAPI00528 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00528(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00528 fpFutureGeneric[string]

func (f manyAPI00528) next(s string) manyAPI00529 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00529(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00529 fpFutureGeneric[string]

func (f manyAPI00529) next(s string) manyAPI00530 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00530(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00530 fpFutureGeneric[string]

func (f manyAPI00530) next(s string) manyAPI00531 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00531(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00531 fpFutureGeneric[string]

func (f manyAPI00531) next(s string) manyAPI00532 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00532(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00532 fpFutureGeneric[string]

func (f manyAPI00532) next(s string) manyAPI00533 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00533(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00533 fpFutureGeneric[string]

func (f manyAPI00533) next(s string) manyAPI00534 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00534(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00534 fpFutureGeneric[string]

func (f manyAPI00534) next(s string) manyAPI00535 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00535(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00535 fpFutureGeneric[string]

func (f manyAPI00535) next(s string) manyAPI00536 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00536(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00536 fpFutureGeneric[string]

func (f manyAPI00536) next(s string) manyAPI00537 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00537(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00537 fpFutureGeneric[string]

func (f manyAPI00537) next(s string) manyAPI00538 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00538(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00538 fpFutureGeneric[string]

func (f manyAPI00538) next(s string) manyAPI00539 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00539(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00539 fpFutureGeneric[string]

func (f manyAPI00539) next(s string) manyAPI00540 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00540(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00540 fpFutureGeneric[string]

func (f manyAPI00540) next(s string) manyAPI00541 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00541(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00541 fpFutureGeneric[string]

func (f manyAPI00541) next(s string) manyAPI00542 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00542(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00542 fpFutureGeneric[string]

func (f manyAPI00542) next(s string) manyAPI00543 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00543(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00543 fpFutureGeneric[string]

func (f manyAPI00543) next(s string) manyAPI00544 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00544(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00544 fpFutureGeneric[string]

func (f manyAPI00544) next(s string) manyAPI00545 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00545(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00545 fpFutureGeneric[string]

func (f manyAPI00545) next(s string) manyAPI00546 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00546(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00546 fpFutureGeneric[string]

func (f manyAPI00546) next(s string) manyAPI00547 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00547(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00547 fpFutureGeneric[string]

func (f manyAPI00547) next(s string) manyAPI00548 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00548(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00548 fpFutureGeneric[string]

func (f manyAPI00548) next(s string) manyAPI00549 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00549(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00549 fpFutureGeneric[string]

func (f manyAPI00549) next(s string) manyAPI00550 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00550(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00550 fpFutureGeneric[string]

func (f manyAPI00550) next(s string) manyAPI00551 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00551(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00551 fpFutureGeneric[string]

func (f manyAPI00551) next(s string) manyAPI00552 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00552(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00552 fpFutureGeneric[string]

func (f manyAPI00552) next(s string) manyAPI00553 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00553(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00553 fpFutureGeneric[string]

func (f manyAPI00553) next(s string) manyAPI00554 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00554(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00554 fpFutureGeneric[string]

func (f manyAPI00554) next(s string) manyAPI00555 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00555(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00555 fpFutureGeneric[string]

func (f manyAPI00555) next(s string) manyAPI00556 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00556(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00556 fpFutureGeneric[string]

func (f manyAPI00556) next(s string) manyAPI00557 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00557(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00557 fpFutureGeneric[string]

func (f manyAPI00557) next(s string) manyAPI00558 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00558(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00558 fpFutureGeneric[string]

func (f manyAPI00558) next(s string) manyAPI00559 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00559(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00559 fpFutureGeneric[string]

func (f manyAPI00559) next(s string) manyAPI00560 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00560(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00560 fpFutureGeneric[string]

func (f manyAPI00560) next(s string) manyAPI00561 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00561(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00561 fpFutureGeneric[string]

func (f manyAPI00561) next(s string) manyAPI00562 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00562(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00562 fpFutureGeneric[string]

func (f manyAPI00562) next(s string) manyAPI00563 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00563(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00563 fpFutureGeneric[string]

func (f manyAPI00563) next(s string) manyAPI00564 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00564(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00564 fpFutureGeneric[string]

func (f manyAPI00564) next(s string) manyAPI00565 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00565(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00565 fpFutureGeneric[string]

func (f manyAPI00565) next(s string) manyAPI00566 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00566(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00566 fpFutureGeneric[string]

func (f manyAPI00566) next(s string) manyAPI00567 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00567(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00567 fpFutureGeneric[string]

func (f manyAPI00567) next(s string) manyAPI00568 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00568(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00568 fpFutureGeneric[string]

func (f manyAPI00568) next(s string) manyAPI00569 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00569(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00569 fpFutureGeneric[string]

func (f manyAPI00569) next(s string) manyAPI00570 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00570(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00570 fpFutureGeneric[string]

func (f manyAPI00570) next(s string) manyAPI00571 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00571(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00571 fpFutureGeneric[string]

func (f manyAPI00571) next(s string) manyAPI00572 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00572(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00572 fpFutureGeneric[string]

func (f manyAPI00572) next(s string) manyAPI00573 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00573(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00573 fpFutureGeneric[string]

func (f manyAPI00573) next(s string) manyAPI00574 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00574(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00574 fpFutureGeneric[string]

func (f manyAPI00574) next(s string) manyAPI00575 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00575(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00575 fpFutureGeneric[string]

func (f manyAPI00575) next(s string) manyAPI00576 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00576(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00576 fpFutureGeneric[string]

func (f manyAPI00576) next(s string) manyAPI00577 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00577(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00577 fpFutureGeneric[string]

func (f manyAPI00577) next(s string) manyAPI00578 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00578(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00578 fpFutureGeneric[string]

func (f manyAPI00578) next(s string) manyAPI00579 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00579(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00579 fpFutureGeneric[string]

func (f manyAPI00579) next(s string) manyAPI00580 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00580(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00580 fpFutureGeneric[string]

func (f manyAPI00580) next(s string) manyAPI00581 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00581(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00581 fpFutureGeneric[string]

func (f manyAPI00581) next(s string) manyAPI00582 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00582(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00582 fpFutureGeneric[string]

func (f manyAPI00582) next(s string) manyAPI00583 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00583(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00583 fpFutureGeneric[string]

func (f manyAPI00583) next(s string) manyAPI00584 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00584(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00584 fpFutureGeneric[string]

func (f manyAPI00584) next(s string) manyAPI00585 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00585(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00585 fpFutureGeneric[string]

func (f manyAPI00585) next(s string) manyAPI00586 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00586(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00586 fpFutureGeneric[string]

func (f manyAPI00586) next(s string) manyAPI00587 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00587(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00587 fpFutureGeneric[string]

func (f manyAPI00587) next(s string) manyAPI00588 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00588(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00588 fpFutureGeneric[string]

func (f manyAPI00588) next(s string) manyAPI00589 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00589(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00589 fpFutureGeneric[string]

func (f manyAPI00589) next(s string) manyAPI00590 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00590(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00590 fpFutureGeneric[string]

func (f manyAPI00590) next(s string) manyAPI00591 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00591(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00591 fpFutureGeneric[string]

func (f manyAPI00591) next(s string) manyAPI00592 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00592(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00592 fpFutureGeneric[string]

func (f manyAPI00592) next(s string) manyAPI00593 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00593(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00593 fpFutureGeneric[string]

func (f manyAPI00593) next(s string) manyAPI00594 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00594(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00594 fpFutureGeneric[string]

func (f manyAPI00594) next(s string) manyAPI00595 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00595(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00595 fpFutureGeneric[string]

func (f manyAPI00595) next(s string) manyAPI00596 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00596(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00596 fpFutureGeneric[string]

func (f manyAPI00596) next(s string) manyAPI00597 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00597(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00597 fpFutureGeneric[string]

func (f manyAPI00597) next(s string) manyAPI00598 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00598(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00598 fpFutureGeneric[string]

func (f manyAPI00598) next(s string) manyAPI00599 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00599(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00599 fpFutureGeneric[string]

func (f manyAPI00599) next(s string) manyAPI00600 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00600(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00600 fpFutureGeneric[string]

func (f manyAPI00600) next(s string) manyAPI00601 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00601(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00601 fpFutureGeneric[string]

func (f manyAPI00601) next(s string) manyAPI00602 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00602(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00602 fpFutureGeneric[string]

func (f manyAPI00602) next(s string) manyAPI00603 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00603(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00603 fpFutureGeneric[string]

func (f manyAPI00603) next(s string) manyAPI00604 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00604(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00604 fpFutureGeneric[string]

func (f manyAPI00604) next(s string) manyAPI00605 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00605(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00605 fpFutureGeneric[string]

func (f manyAPI00605) next(s string) manyAPI00606 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00606(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00606 fpFutureGeneric[string]

func (f manyAPI00606) next(s string) manyAPI00607 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00607(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00607 fpFutureGeneric[string]

func (f manyAPI00607) next(s string) manyAPI00608 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00608(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00608 fpFutureGeneric[string]

func (f manyAPI00608) next(s string) manyAPI00609 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00609(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00609 fpFutureGeneric[string]

func (f manyAPI00609) next(s string) manyAPI00610 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00610(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00610 fpFutureGeneric[string]

func (f manyAPI00610) next(s string) manyAPI00611 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00611(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00611 fpFutureGeneric[string]

func (f manyAPI00611) next(s string) manyAPI00612 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00612(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00612 fpFutureGeneric[string]

func (f manyAPI00612) next(s string) manyAPI00613 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00613(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00613 fpFutureGeneric[string]

func (f manyAPI00613) next(s string) manyAPI00614 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00614(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00614 fpFutureGeneric[string]

func (f manyAPI00614) next(s string) manyAPI00615 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00615(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00615 fpFutureGeneric[string]

func (f manyAPI00615) next(s string) manyAPI00616 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00616(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00616 fpFutureGeneric[string]

func (f manyAPI00616) next(s string) manyAPI00617 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00617(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00617 fpFutureGeneric[string]

func (f manyAPI00617) next(s string) manyAPI00618 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00618(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00618 fpFutureGeneric[string]

func (f manyAPI00618) next(s string) manyAPI00619 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00619(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00619 fpFutureGeneric[string]

func (f manyAPI00619) next(s string) manyAPI00620 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00620(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00620 fpFutureGeneric[string]

func (f manyAPI00620) next(s string) manyAPI00621 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00621(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00621 fpFutureGeneric[string]

func (f manyAPI00621) next(s string) manyAPI00622 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00622(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00622 fpFutureGeneric[string]

func (f manyAPI00622) next(s string) manyAPI00623 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00623(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00623 fpFutureGeneric[string]

func (f manyAPI00623) next(s string) manyAPI00624 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00624(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00624 fpFutureGeneric[string]

func (f manyAPI00624) next(s string) manyAPI00625 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00625(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00625 fpFutureGeneric[string]

func (f manyAPI00625) next(s string) manyAPI00626 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00626(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00626 fpFutureGeneric[string]

func (f manyAPI00626) next(s string) manyAPI00627 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00627(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00627 fpFutureGeneric[string]

func (f manyAPI00627) next(s string) manyAPI00628 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00628(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00628 fpFutureGeneric[string]

func (f manyAPI00628) next(s string) manyAPI00629 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00629(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00629 fpFutureGeneric[string]

func (f manyAPI00629) next(s string) manyAPI00630 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00630(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00630 fpFutureGeneric[string]

func (f manyAPI00630) next(s string) manyAPI00631 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00631(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00631 fpFutureGeneric[string]

func (f manyAPI00631) next(s string) manyAPI00632 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00632(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00632 fpFutureGeneric[string]

func (f manyAPI00632) next(s string) manyAPI00633 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00633(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00633 fpFutureGeneric[string]

func (f manyAPI00633) next(s string) manyAPI00634 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00634(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00634 fpFutureGeneric[string]

func (f manyAPI00634) next(s string) manyAPI00635 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00635(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00635 fpFutureGeneric[string]

func (f manyAPI00635) next(s string) manyAPI00636 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00636(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00636 fpFutureGeneric[string]

func (f manyAPI00636) next(s string) manyAPI00637 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00637(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00637 fpFutureGeneric[string]

func (f manyAPI00637) next(s string) manyAPI00638 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00638(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00638 fpFutureGeneric[string]

func (f manyAPI00638) next(s string) manyAPI00639 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00639(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00639 fpFutureGeneric[string]

func (f manyAPI00639) next(s string) manyAPI00640 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00640(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00640 fpFutureGeneric[string]

func (f manyAPI00640) next(s string) manyAPI00641 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00641(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00641 fpFutureGeneric[string]

func (f manyAPI00641) next(s string) manyAPI00642 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00642(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00642 fpFutureGeneric[string]

func (f manyAPI00642) next(s string) manyAPI00643 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00643(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00643 fpFutureGeneric[string]

func (f manyAPI00643) next(s string) manyAPI00644 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00644(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00644 fpFutureGeneric[string]

func (f manyAPI00644) next(s string) manyAPI00645 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00645(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00645 fpFutureGeneric[string]

func (f manyAPI00645) next(s string) manyAPI00646 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00646(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00646 fpFutureGeneric[string]

func (f manyAPI00646) next(s string) manyAPI00647 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00647(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00647 fpFutureGeneric[string]

func (f manyAPI00647) next(s string) manyAPI00648 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00648(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00648 fpFutureGeneric[string]

func (f manyAPI00648) next(s string) manyAPI00649 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00649(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00649 fpFutureGeneric[string]

func (f manyAPI00649) next(s string) manyAPI00650 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00650(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00650 fpFutureGeneric[string]

func (f manyAPI00650) next(s string) manyAPI00651 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00651(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00651 fpFutureGeneric[string]

func (f manyAPI00651) next(s string) manyAPI00652 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00652(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00652 fpFutureGeneric[string]

func (f manyAPI00652) next(s string) manyAPI00653 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00653(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00653 fpFutureGeneric[string]

func (f manyAPI00653) next(s string) manyAPI00654 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00654(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00654 fpFutureGeneric[string]

func (f manyAPI00654) next(s string) manyAPI00655 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00655(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00655 fpFutureGeneric[string]

func (f manyAPI00655) next(s string) manyAPI00656 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00656(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00656 fpFutureGeneric[string]

func (f manyAPI00656) next(s string) manyAPI00657 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00657(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00657 fpFutureGeneric[string]

func (f manyAPI00657) next(s string) manyAPI00658 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00658(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00658 fpFutureGeneric[string]

func (f manyAPI00658) next(s string) manyAPI00659 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00659(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00659 fpFutureGeneric[string]

func (f manyAPI00659) next(s string) manyAPI00660 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00660(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00660 fpFutureGeneric[string]

func (f manyAPI00660) next(s string) manyAPI00661 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00661(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00661 fpFutureGeneric[string]

func (f manyAPI00661) next(s string) manyAPI00662 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00662(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00662 fpFutureGeneric[string]

func (f manyAPI00662) next(s string) manyAPI00663 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00663(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00663 fpFutureGeneric[string]

func (f manyAPI00663) next(s string) manyAPI00664 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00664(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00664 fpFutureGeneric[string]

func (f manyAPI00664) next(s string) manyAPI00665 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00665(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00665 fpFutureGeneric[string]

func (f manyAPI00665) next(s string) manyAPI00666 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00666(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00666 fpFutureGeneric[string]

func (f manyAPI00666) next(s string) manyAPI00667 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00667(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00667 fpFutureGeneric[string]

func (f manyAPI00667) next(s string) manyAPI00668 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00668(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00668 fpFutureGeneric[string]

func (f manyAPI00668) next(s string) manyAPI00669 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00669(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00669 fpFutureGeneric[string]

func (f manyAPI00669) next(s string) manyAPI00670 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00670(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00670 fpFutureGeneric[string]

func (f manyAPI00670) next(s string) manyAPI00671 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00671(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00671 fpFutureGeneric[string]

func (f manyAPI00671) next(s string) manyAPI00672 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00672(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00672 fpFutureGeneric[string]

func (f manyAPI00672) next(s string) manyAPI00673 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00673(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00673 fpFutureGeneric[string]

func (f manyAPI00673) next(s string) manyAPI00674 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00674(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00674 fpFutureGeneric[string]

func (f manyAPI00674) next(s string) manyAPI00675 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00675(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00675 fpFutureGeneric[string]

func (f manyAPI00675) next(s string) manyAPI00676 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00676(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00676 fpFutureGeneric[string]

func (f manyAPI00676) next(s string) manyAPI00677 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00677(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00677 fpFutureGeneric[string]

func (f manyAPI00677) next(s string) manyAPI00678 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00678(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00678 fpFutureGeneric[string]

func (f manyAPI00678) next(s string) manyAPI00679 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00679(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00679 fpFutureGeneric[string]

func (f manyAPI00679) next(s string) manyAPI00680 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00680(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00680 fpFutureGeneric[string]

func (f manyAPI00680) next(s string) manyAPI00681 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00681(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00681 fpFutureGeneric[string]

func (f manyAPI00681) next(s string) manyAPI00682 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00682(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00682 fpFutureGeneric[string]

func (f manyAPI00682) next(s string) manyAPI00683 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00683(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00683 fpFutureGeneric[string]

func (f manyAPI00683) next(s string) manyAPI00684 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00684(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00684 fpFutureGeneric[string]

func (f manyAPI00684) next(s string) manyAPI00685 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00685(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00685 fpFutureGeneric[string]

func (f manyAPI00685) next(s string) manyAPI00686 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00686(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00686 fpFutureGeneric[string]

func (f manyAPI00686) next(s string) manyAPI00687 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00687(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00687 fpFutureGeneric[string]

func (f manyAPI00687) next(s string) manyAPI00688 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00688(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00688 fpFutureGeneric[string]

func (f manyAPI00688) next(s string) manyAPI00689 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00689(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00689 fpFutureGeneric[string]

func (f manyAPI00689) next(s string) manyAPI00690 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00690(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00690 fpFutureGeneric[string]

func (f manyAPI00690) next(s string) manyAPI00691 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00691(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00691 fpFutureGeneric[string]

func (f manyAPI00691) next(s string) manyAPI00692 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00692(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00692 fpFutureGeneric[string]

func (f manyAPI00692) next(s string) manyAPI00693 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00693(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00693 fpFutureGeneric[string]

func (f manyAPI00693) next(s string) manyAPI00694 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00694(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00694 fpFutureGeneric[string]

func (f manyAPI00694) next(s string) manyAPI00695 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00695(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00695 fpFutureGeneric[string]

func (f manyAPI00695) next(s string) manyAPI00696 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00696(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00696 fpFutureGeneric[string]

func (f manyAPI00696) next(s string) manyAPI00697 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00697(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00697 fpFutureGeneric[string]

func (f manyAPI00697) next(s string) manyAPI00698 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00698(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00698 fpFutureGeneric[string]

func (f manyAPI00698) next(s string) manyAPI00699 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00699(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00699 fpFutureGeneric[string]

func (f manyAPI00699) next(s string) manyAPI00700 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00700(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00700 fpFutureGeneric[string]

func (f manyAPI00700) next(s string) manyAPI00701 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00701(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00701 fpFutureGeneric[string]

func (f manyAPI00701) next(s string) manyAPI00702 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00702(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00702 fpFutureGeneric[string]

func (f manyAPI00702) next(s string) manyAPI00703 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00703(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00703 fpFutureGeneric[string]

func (f manyAPI00703) next(s string) manyAPI00704 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00704(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00704 fpFutureGeneric[string]

func (f manyAPI00704) next(s string) manyAPI00705 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00705(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00705 fpFutureGeneric[string]

func (f manyAPI00705) next(s string) manyAPI00706 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00706(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00706 fpFutureGeneric[string]

func (f manyAPI00706) next(s string) manyAPI00707 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00707(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00707 fpFutureGeneric[string]

func (f manyAPI00707) next(s string) manyAPI00708 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00708(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00708 fpFutureGeneric[string]

func (f manyAPI00708) next(s string) manyAPI00709 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00709(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00709 fpFutureGeneric[string]

func (f manyAPI00709) next(s string) manyAPI00710 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00710(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00710 fpFutureGeneric[string]

func (f manyAPI00710) next(s string) manyAPI00711 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00711(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00711 fpFutureGeneric[string]

func (f manyAPI00711) next(s string) manyAPI00712 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00712(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00712 fpFutureGeneric[string]

func (f manyAPI00712) next(s string) manyAPI00713 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00713(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00713 fpFutureGeneric[string]

func (f manyAPI00713) next(s string) manyAPI00714 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00714(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00714 fpFutureGeneric[string]

func (f manyAPI00714) next(s string) manyAPI00715 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00715(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00715 fpFutureGeneric[string]

func (f manyAPI00715) next(s string) manyAPI00716 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00716(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00716 fpFutureGeneric[string]

func (f manyAPI00716) next(s string) manyAPI00717 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00717(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00717 fpFutureGeneric[string]

func (f manyAPI00717) next(s string) manyAPI00718 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00718(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00718 fpFutureGeneric[string]

func (f manyAPI00718) next(s string) manyAPI00719 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00719(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00719 fpFutureGeneric[string]

func (f manyAPI00719) next(s string) manyAPI00720 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00720(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00720 fpFutureGeneric[string]

func (f manyAPI00720) next(s string) manyAPI00721 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00721(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00721 fpFutureGeneric[string]

func (f manyAPI00721) next(s string) manyAPI00722 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00722(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00722 fpFutureGeneric[string]

func (f manyAPI00722) next(s string) manyAPI00723 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00723(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00723 fpFutureGeneric[string]

func (f manyAPI00723) next(s string) manyAPI00724 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00724(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00724 fpFutureGeneric[string]

func (f manyAPI00724) next(s string) manyAPI00725 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00725(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00725 fpFutureGeneric[string]

func (f manyAPI00725) next(s string) manyAPI00726 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00726(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00726 fpFutureGeneric[string]

func (f manyAPI00726) next(s string) manyAPI00727 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00727(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00727 fpFutureGeneric[string]

func (f manyAPI00727) next(s string) manyAPI00728 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00728(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00728 fpFutureGeneric[string]

func (f manyAPI00728) next(s string) manyAPI00729 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00729(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00729 fpFutureGeneric[string]

func (f manyAPI00729) next(s string) manyAPI00730 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00730(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00730 fpFutureGeneric[string]

func (f manyAPI00730) next(s string) manyAPI00731 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00731(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00731 fpFutureGeneric[string]

func (f manyAPI00731) next(s string) manyAPI00732 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00732(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00732 fpFutureGeneric[string]

func (f manyAPI00732) next(s string) manyAPI00733 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00733(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00733 fpFutureGeneric[string]

func (f manyAPI00733) next(s string) manyAPI00734 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00734(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00734 fpFutureGeneric[string]

func (f manyAPI00734) next(s string) manyAPI00735 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00735(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00735 fpFutureGeneric[string]

func (f manyAPI00735) next(s string) manyAPI00736 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00736(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00736 fpFutureGeneric[string]

func (f manyAPI00736) next(s string) manyAPI00737 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00737(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00737 fpFutureGeneric[string]

func (f manyAPI00737) next(s string) manyAPI00738 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00738(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00738 fpFutureGeneric[string]

func (f manyAPI00738) next(s string) manyAPI00739 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00739(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00739 fpFutureGeneric[string]

func (f manyAPI00739) next(s string) manyAPI00740 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00740(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00740 fpFutureGeneric[string]

func (f manyAPI00740) next(s string) manyAPI00741 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00741(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00741 fpFutureGeneric[string]

func (f manyAPI00741) next(s string) manyAPI00742 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00742(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00742 fpFutureGeneric[string]

func (f manyAPI00742) next(s string) manyAPI00743 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00743(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00743 fpFutureGeneric[string]

func (f manyAPI00743) next(s string) manyAPI00744 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00744(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00744 fpFutureGeneric[string]

func (f manyAPI00744) next(s string) manyAPI00745 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00745(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00745 fpFutureGeneric[string]

func (f manyAPI00745) next(s string) manyAPI00746 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00746(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00746 fpFutureGeneric[string]

func (f manyAPI00746) next(s string) manyAPI00747 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00747(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00747 fpFutureGeneric[string]

func (f manyAPI00747) next(s string) manyAPI00748 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00748(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00748 fpFutureGeneric[string]

func (f manyAPI00748) next(s string) manyAPI00749 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00749(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00749 fpFutureGeneric[string]

func (f manyAPI00749) next(s string) manyAPI00750 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00750(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00750 fpFutureGeneric[string]

func (f manyAPI00750) next(s string) manyAPI00751 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00751(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00751 fpFutureGeneric[string]

func (f manyAPI00751) next(s string) manyAPI00752 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00752(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00752 fpFutureGeneric[string]

func (f manyAPI00752) next(s string) manyAPI00753 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00753(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00753 fpFutureGeneric[string]

func (f manyAPI00753) next(s string) manyAPI00754 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00754(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00754 fpFutureGeneric[string]

func (f manyAPI00754) next(s string) manyAPI00755 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00755(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00755 fpFutureGeneric[string]

func (f manyAPI00755) next(s string) manyAPI00756 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00756(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00756 fpFutureGeneric[string]

func (f manyAPI00756) next(s string) manyAPI00757 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00757(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00757 fpFutureGeneric[string]

func (f manyAPI00757) next(s string) manyAPI00758 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00758(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00758 fpFutureGeneric[string]

func (f manyAPI00758) next(s string) manyAPI00759 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00759(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00759 fpFutureGeneric[string]

func (f manyAPI00759) next(s string) manyAPI00760 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00760(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00760 fpFutureGeneric[string]

func (f manyAPI00760) next(s string) manyAPI00761 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00761(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00761 fpFutureGeneric[string]

func (f manyAPI00761) next(s string) manyAPI00762 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00762(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00762 fpFutureGeneric[string]

func (f manyAPI00762) next(s string) manyAPI00763 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00763(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00763 fpFutureGeneric[string]

func (f manyAPI00763) next(s string) manyAPI00764 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00764(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00764 fpFutureGeneric[string]

func (f manyAPI00764) next(s string) manyAPI00765 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00765(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00765 fpFutureGeneric[string]

func (f manyAPI00765) next(s string) manyAPI00766 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00766(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00766 fpFutureGeneric[string]

func (f manyAPI00766) next(s string) manyAPI00767 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00767(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00767 fpFutureGeneric[string]

func (f manyAPI00767) next(s string) manyAPI00768 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00768(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00768 fpFutureGeneric[string]

func (f manyAPI00768) next(s string) manyAPI00769 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00769(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00769 fpFutureGeneric[string]

func (f manyAPI00769) next(s string) manyAPI00770 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00770(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00770 fpFutureGeneric[string]

func (f manyAPI00770) next(s string) manyAPI00771 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00771(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00771 fpFutureGeneric[string]

func (f manyAPI00771) next(s string) manyAPI00772 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00772(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00772 fpFutureGeneric[string]

func (f manyAPI00772) next(s string) manyAPI00773 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00773(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00773 fpFutureGeneric[string]

func (f manyAPI00773) next(s string) manyAPI00774 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00774(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00774 fpFutureGeneric[string]

func (f manyAPI00774) next(s string) manyAPI00775 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00775(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00775 fpFutureGeneric[string]

func (f manyAPI00775) next(s string) manyAPI00776 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00776(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00776 fpFutureGeneric[string]

func (f manyAPI00776) next(s string) manyAPI00777 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00777(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00777 fpFutureGeneric[string]

func (f manyAPI00777) next(s string) manyAPI00778 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00778(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00778 fpFutureGeneric[string]

func (f manyAPI00778) next(s string) manyAPI00779 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00779(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00779 fpFutureGeneric[string]

func (f manyAPI00779) next(s string) manyAPI00780 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00780(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00780 fpFutureGeneric[string]

func (f manyAPI00780) next(s string) manyAPI00781 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00781(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00781 fpFutureGeneric[string]

func (f manyAPI00781) next(s string) manyAPI00782 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00782(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00782 fpFutureGeneric[string]

func (f manyAPI00782) next(s string) manyAPI00783 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00783(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00783 fpFutureGeneric[string]

func (f manyAPI00783) next(s string) manyAPI00784 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00784(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00784 fpFutureGeneric[string]

func (f manyAPI00784) next(s string) manyAPI00785 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00785(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00785 fpFutureGeneric[string]

func (f manyAPI00785) next(s string) manyAPI00786 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00786(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00786 fpFutureGeneric[string]

func (f manyAPI00786) next(s string) manyAPI00787 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00787(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00787 fpFutureGeneric[string]

func (f manyAPI00787) next(s string) manyAPI00788 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00788(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00788 fpFutureGeneric[string]

func (f manyAPI00788) next(s string) manyAPI00789 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00789(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00789 fpFutureGeneric[string]

func (f manyAPI00789) next(s string) manyAPI00790 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00790(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00790 fpFutureGeneric[string]

func (f manyAPI00790) next(s string) manyAPI00791 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00791(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00791 fpFutureGeneric[string]

func (f manyAPI00791) next(s string) manyAPI00792 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00792(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00792 fpFutureGeneric[string]

func (f manyAPI00792) next(s string) manyAPI00793 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00793(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00793 fpFutureGeneric[string]

func (f manyAPI00793) next(s string) manyAPI00794 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00794(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00794 fpFutureGeneric[string]

func (f manyAPI00794) next(s string) manyAPI00795 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00795(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00795 fpFutureGeneric[string]

func (f manyAPI00795) next(s string) manyAPI00796 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00796(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00796 fpFutureGeneric[string]

func (f manyAPI00796) next(s string) manyAPI00797 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00797(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00797 fpFutureGeneric[string]

func (f manyAPI00797) next(s string) manyAPI00798 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00798(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00798 fpFutureGeneric[string]

func (f manyAPI00798) next(s string) manyAPI00799 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00799(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00799 fpFutureGeneric[string]

func (f manyAPI00799) next(s string) manyAPI00800 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00800(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00800 fpFutureGeneric[string]

func (f manyAPI00800) next(s string) manyAPI00801 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00801(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00801 fpFutureGeneric[string]

func (f manyAPI00801) next(s string) manyAPI00802 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00802(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00802 fpFutureGeneric[string]

func (f manyAPI00802) next(s string) manyAPI00803 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00803(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00803 fpFutureGeneric[string]

func (f manyAPI00803) next(s string) manyAPI00804 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00804(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00804 fpFutureGeneric[string]

func (f manyAPI00804) next(s string) manyAPI00805 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00805(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00805 fpFutureGeneric[string]

func (f manyAPI00805) next(s string) manyAPI00806 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00806(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00806 fpFutureGeneric[string]

func (f manyAPI00806) next(s string) manyAPI00807 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00807(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00807 fpFutureGeneric[string]

func (f manyAPI00807) next(s string) manyAPI00808 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00808(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00808 fpFutureGeneric[string]

func (f manyAPI00808) next(s string) manyAPI00809 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00809(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00809 fpFutureGeneric[string]

func (f manyAPI00809) next(s string) manyAPI00810 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00810(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00810 fpFutureGeneric[string]

func (f manyAPI00810) next(s string) manyAPI00811 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00811(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00811 fpFutureGeneric[string]

func (f manyAPI00811) next(s string) manyAPI00812 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00812(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00812 fpFutureGeneric[string]

func (f manyAPI00812) next(s string) manyAPI00813 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00813(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00813 fpFutureGeneric[string]

func (f manyAPI00813) next(s string) manyAPI00814 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00814(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00814 fpFutureGeneric[string]

func (f manyAPI00814) next(s string) manyAPI00815 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00815(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00815 fpFutureGeneric[string]

func (f manyAPI00815) next(s string) manyAPI00816 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00816(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00816 fpFutureGeneric[string]

func (f manyAPI00816) next(s string) manyAPI00817 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00817(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00817 fpFutureGeneric[string]

func (f manyAPI00817) next(s string) manyAPI00818 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00818(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00818 fpFutureGeneric[string]

func (f manyAPI00818) next(s string) manyAPI00819 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00819(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00819 fpFutureGeneric[string]

func (f manyAPI00819) next(s string) manyAPI00820 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00820(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00820 fpFutureGeneric[string]

func (f manyAPI00820) next(s string) manyAPI00821 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00821(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00821 fpFutureGeneric[string]

func (f manyAPI00821) next(s string) manyAPI00822 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00822(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00822 fpFutureGeneric[string]

func (f manyAPI00822) next(s string) manyAPI00823 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00823(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00823 fpFutureGeneric[string]

func (f manyAPI00823) next(s string) manyAPI00824 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00824(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00824 fpFutureGeneric[string]

func (f manyAPI00824) next(s string) manyAPI00825 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00825(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00825 fpFutureGeneric[string]

func (f manyAPI00825) next(s string) manyAPI00826 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00826(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00826 fpFutureGeneric[string]

func (f manyAPI00826) next(s string) manyAPI00827 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00827(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00827 fpFutureGeneric[string]

func (f manyAPI00827) next(s string) manyAPI00828 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00828(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00828 fpFutureGeneric[string]

func (f manyAPI00828) next(s string) manyAPI00829 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00829(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00829 fpFutureGeneric[string]

func (f manyAPI00829) next(s string) manyAPI00830 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00830(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00830 fpFutureGeneric[string]

func (f manyAPI00830) next(s string) manyAPI00831 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00831(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00831 fpFutureGeneric[string]

func (f manyAPI00831) next(s string) manyAPI00832 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00832(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00832 fpFutureGeneric[string]

func (f manyAPI00832) next(s string) manyAPI00833 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00833(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00833 fpFutureGeneric[string]

func (f manyAPI00833) next(s string) manyAPI00834 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00834(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00834 fpFutureGeneric[string]

func (f manyAPI00834) next(s string) manyAPI00835 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00835(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00835 fpFutureGeneric[string]

func (f manyAPI00835) next(s string) manyAPI00836 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00836(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00836 fpFutureGeneric[string]

func (f manyAPI00836) next(s string) manyAPI00837 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00837(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00837 fpFutureGeneric[string]

func (f manyAPI00837) next(s string) manyAPI00838 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00838(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00838 fpFutureGeneric[string]

func (f manyAPI00838) next(s string) manyAPI00839 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00839(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00839 fpFutureGeneric[string]

func (f manyAPI00839) next(s string) manyAPI00840 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00840(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00840 fpFutureGeneric[string]

func (f manyAPI00840) next(s string) manyAPI00841 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00841(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00841 fpFutureGeneric[string]

func (f manyAPI00841) next(s string) manyAPI00842 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00842(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00842 fpFutureGeneric[string]

func (f manyAPI00842) next(s string) manyAPI00843 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00843(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00843 fpFutureGeneric[string]

func (f manyAPI00843) next(s string) manyAPI00844 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00844(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00844 fpFutureGeneric[string]

func (f manyAPI00844) next(s string) manyAPI00845 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00845(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00845 fpFutureGeneric[string]

func (f manyAPI00845) next(s string) manyAPI00846 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00846(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00846 fpFutureGeneric[string]

func (f manyAPI00846) next(s string) manyAPI00847 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00847(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00847 fpFutureGeneric[string]

func (f manyAPI00847) next(s string) manyAPI00848 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00848(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00848 fpFutureGeneric[string]

func (f manyAPI00848) next(s string) manyAPI00849 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00849(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00849 fpFutureGeneric[string]

func (f manyAPI00849) next(s string) manyAPI00850 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00850(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00850 fpFutureGeneric[string]

func (f manyAPI00850) next(s string) manyAPI00851 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00851(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00851 fpFutureGeneric[string]

func (f manyAPI00851) next(s string) manyAPI00852 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00852(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00852 fpFutureGeneric[string]

func (f manyAPI00852) next(s string) manyAPI00853 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00853(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00853 fpFutureGeneric[string]

func (f manyAPI00853) next(s string) manyAPI00854 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00854(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00854 fpFutureGeneric[string]

func (f manyAPI00854) next(s string) manyAPI00855 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00855(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00855 fpFutureGeneric[string]

func (f manyAPI00855) next(s string) manyAPI00856 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00856(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00856 fpFutureGeneric[string]

func (f manyAPI00856) next(s string) manyAPI00857 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00857(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00857 fpFutureGeneric[string]

func (f manyAPI00857) next(s string) manyAPI00858 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00858(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00858 fpFutureGeneric[string]

func (f manyAPI00858) next(s string) manyAPI00859 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00859(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00859 fpFutureGeneric[string]

func (f manyAPI00859) next(s string) manyAPI00860 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00860(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00860 fpFutureGeneric[string]

func (f manyAPI00860) next(s string) manyAPI00861 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00861(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00861 fpFutureGeneric[string]

func (f manyAPI00861) next(s string) manyAPI00862 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00862(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00862 fpFutureGeneric[string]

func (f manyAPI00862) next(s string) manyAPI00863 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00863(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00863 fpFutureGeneric[string]

func (f manyAPI00863) next(s string) manyAPI00864 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00864(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00864 fpFutureGeneric[string]

func (f manyAPI00864) next(s string) manyAPI00865 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00865(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00865 fpFutureGeneric[string]

func (f manyAPI00865) next(s string) manyAPI00866 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00866(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00866 fpFutureGeneric[string]

func (f manyAPI00866) next(s string) manyAPI00867 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00867(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00867 fpFutureGeneric[string]

func (f manyAPI00867) next(s string) manyAPI00868 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00868(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00868 fpFutureGeneric[string]

func (f manyAPI00868) next(s string) manyAPI00869 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00869(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00869 fpFutureGeneric[string]

func (f manyAPI00869) next(s string) manyAPI00870 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00870(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00870 fpFutureGeneric[string]

func (f manyAPI00870) next(s string) manyAPI00871 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00871(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00871 fpFutureGeneric[string]

func (f manyAPI00871) next(s string) manyAPI00872 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00872(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00872 fpFutureGeneric[string]

func (f manyAPI00872) next(s string) manyAPI00873 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00873(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00873 fpFutureGeneric[string]

func (f manyAPI00873) next(s string) manyAPI00874 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00874(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00874 fpFutureGeneric[string]

func (f manyAPI00874) next(s string) manyAPI00875 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00875(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00875 fpFutureGeneric[string]

func (f manyAPI00875) next(s string) manyAPI00876 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00876(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00876 fpFutureGeneric[string]

func (f manyAPI00876) next(s string) manyAPI00877 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00877(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00877 fpFutureGeneric[string]

func (f manyAPI00877) next(s string) manyAPI00878 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00878(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00878 fpFutureGeneric[string]

func (f manyAPI00878) next(s string) manyAPI00879 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00879(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00879 fpFutureGeneric[string]

func (f manyAPI00879) next(s string) manyAPI00880 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00880(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00880 fpFutureGeneric[string]

func (f manyAPI00880) next(s string) manyAPI00881 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00881(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00881 fpFutureGeneric[string]

func (f manyAPI00881) next(s string) manyAPI00882 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00882(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00882 fpFutureGeneric[string]

func (f manyAPI00882) next(s string) manyAPI00883 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00883(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00883 fpFutureGeneric[string]

func (f manyAPI00883) next(s string) manyAPI00884 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00884(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00884 fpFutureGeneric[string]

func (f manyAPI00884) next(s string) manyAPI00885 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00885(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00885 fpFutureGeneric[string]

func (f manyAPI00885) next(s string) manyAPI00886 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00886(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00886 fpFutureGeneric[string]

func (f manyAPI00886) next(s string) manyAPI00887 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00887(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00887 fpFutureGeneric[string]

func (f manyAPI00887) next(s string) manyAPI00888 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00888(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00888 fpFutureGeneric[string]

func (f manyAPI00888) next(s string) manyAPI00889 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00889(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00889 fpFutureGeneric[string]

func (f manyAPI00889) next(s string) manyAPI00890 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00890(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00890 fpFutureGeneric[string]

func (f manyAPI00890) next(s string) manyAPI00891 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00891(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00891 fpFutureGeneric[string]

func (f manyAPI00891) next(s string) manyAPI00892 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00892(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00892 fpFutureGeneric[string]

func (f manyAPI00892) next(s string) manyAPI00893 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00893(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00893 fpFutureGeneric[string]

func (f manyAPI00893) next(s string) manyAPI00894 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00894(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00894 fpFutureGeneric[string]

func (f manyAPI00894) next(s string) manyAPI00895 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00895(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00895 fpFutureGeneric[string]

func (f manyAPI00895) next(s string) manyAPI00896 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00896(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00896 fpFutureGeneric[string]

func (f manyAPI00896) next(s string) manyAPI00897 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00897(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00897 fpFutureGeneric[string]

func (f manyAPI00897) next(s string) manyAPI00898 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00898(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00898 fpFutureGeneric[string]

func (f manyAPI00898) next(s string) manyAPI00899 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00899(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00899 fpFutureGeneric[string]

func (f manyAPI00899) next(s string) manyAPI00900 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00900(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00900 fpFutureGeneric[string]

func (f manyAPI00900) next(s string) manyAPI00901 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00901(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00901 fpFutureGeneric[string]

func (f manyAPI00901) next(s string) manyAPI00902 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00902(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00902 fpFutureGeneric[string]

func (f manyAPI00902) next(s string) manyAPI00903 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00903(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00903 fpFutureGeneric[string]

func (f manyAPI00903) next(s string) manyAPI00904 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00904(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00904 fpFutureGeneric[string]

func (f manyAPI00904) next(s string) manyAPI00905 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00905(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00905 fpFutureGeneric[string]

func (f manyAPI00905) next(s string) manyAPI00906 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00906(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00906 fpFutureGeneric[string]

func (f manyAPI00906) next(s string) manyAPI00907 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00907(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00907 fpFutureGeneric[string]

func (f manyAPI00907) next(s string) manyAPI00908 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00908(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00908 fpFutureGeneric[string]

func (f manyAPI00908) next(s string) manyAPI00909 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00909(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00909 fpFutureGeneric[string]

func (f manyAPI00909) next(s string) manyAPI00910 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00910(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00910 fpFutureGeneric[string]

func (f manyAPI00910) next(s string) manyAPI00911 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00911(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00911 fpFutureGeneric[string]

func (f manyAPI00911) next(s string) manyAPI00912 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00912(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00912 fpFutureGeneric[string]

func (f manyAPI00912) next(s string) manyAPI00913 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00913(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00913 fpFutureGeneric[string]

func (f manyAPI00913) next(s string) manyAPI00914 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00914(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00914 fpFutureGeneric[string]

func (f manyAPI00914) next(s string) manyAPI00915 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00915(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00915 fpFutureGeneric[string]

func (f manyAPI00915) next(s string) manyAPI00916 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00916(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00916 fpFutureGeneric[string]

func (f manyAPI00916) next(s string) manyAPI00917 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00917(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00917 fpFutureGeneric[string]

func (f manyAPI00917) next(s string) manyAPI00918 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00918(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00918 fpFutureGeneric[string]

func (f manyAPI00918) next(s string) manyAPI00919 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00919(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00919 fpFutureGeneric[string]

func (f manyAPI00919) next(s string) manyAPI00920 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00920(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00920 fpFutureGeneric[string]

func (f manyAPI00920) next(s string) manyAPI00921 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00921(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00921 fpFutureGeneric[string]

func (f manyAPI00921) next(s string) manyAPI00922 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00922(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00922 fpFutureGeneric[string]

func (f manyAPI00922) next(s string) manyAPI00923 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00923(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00923 fpFutureGeneric[string]

func (f manyAPI00923) next(s string) manyAPI00924 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00924(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00924 fpFutureGeneric[string]

func (f manyAPI00924) next(s string) manyAPI00925 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00925(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00925 fpFutureGeneric[string]

func (f manyAPI00925) next(s string) manyAPI00926 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00926(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00926 fpFutureGeneric[string]

func (f manyAPI00926) next(s string) manyAPI00927 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00927(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00927 fpFutureGeneric[string]

func (f manyAPI00927) next(s string) manyAPI00928 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00928(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00928 fpFutureGeneric[string]

func (f manyAPI00928) next(s string) manyAPI00929 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00929(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00929 fpFutureGeneric[string]

func (f manyAPI00929) next(s string) manyAPI00930 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00930(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00930 fpFutureGeneric[string]

func (f manyAPI00930) next(s string) manyAPI00931 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00931(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00931 fpFutureGeneric[string]

func (f manyAPI00931) next(s string) manyAPI00932 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00932(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00932 fpFutureGeneric[string]

func (f manyAPI00932) next(s string) manyAPI00933 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00933(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00933 fpFutureGeneric[string]

func (f manyAPI00933) next(s string) manyAPI00934 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00934(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00934 fpFutureGeneric[string]

func (f manyAPI00934) next(s string) manyAPI00935 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00935(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00935 fpFutureGeneric[string]

func (f manyAPI00935) next(s string) manyAPI00936 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00936(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00936 fpFutureGeneric[string]

func (f manyAPI00936) next(s string) manyAPI00937 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00937(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00937 fpFutureGeneric[string]

func (f manyAPI00937) next(s string) manyAPI00938 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00938(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00938 fpFutureGeneric[string]

func (f manyAPI00938) next(s string) manyAPI00939 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00939(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00939 fpFutureGeneric[string]

func (f manyAPI00939) next(s string) manyAPI00940 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00940(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00940 fpFutureGeneric[string]

func (f manyAPI00940) next(s string) manyAPI00941 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00941(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00941 fpFutureGeneric[string]

func (f manyAPI00941) next(s string) manyAPI00942 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00942(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00942 fpFutureGeneric[string]

func (f manyAPI00942) next(s string) manyAPI00943 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00943(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00943 fpFutureGeneric[string]

func (f manyAPI00943) next(s string) manyAPI00944 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00944(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00944 fpFutureGeneric[string]

func (f manyAPI00944) next(s string) manyAPI00945 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00945(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00945 fpFutureGeneric[string]

func (f manyAPI00945) next(s string) manyAPI00946 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00946(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00946 fpFutureGeneric[string]

func (f manyAPI00946) next(s string) manyAPI00947 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00947(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00947 fpFutureGeneric[string]

func (f manyAPI00947) next(s string) manyAPI00948 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00948(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00948 fpFutureGeneric[string]

func (f manyAPI00948) next(s string) manyAPI00949 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00949(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00949 fpFutureGeneric[string]

func (f manyAPI00949) next(s string) manyAPI00950 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00950(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00950 fpFutureGeneric[string]

func (f manyAPI00950) next(s string) manyAPI00951 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00951(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00951 fpFutureGeneric[string]

func (f manyAPI00951) next(s string) manyAPI00952 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00952(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00952 fpFutureGeneric[string]

func (f manyAPI00952) next(s string) manyAPI00953 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00953(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00953 fpFutureGeneric[string]

func (f manyAPI00953) next(s string) manyAPI00954 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00954(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00954 fpFutureGeneric[string]

func (f manyAPI00954) next(s string) manyAPI00955 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00955(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00955 fpFutureGeneric[string]

func (f manyAPI00955) next(s string) manyAPI00956 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00956(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00956 fpFutureGeneric[string]

func (f manyAPI00956) next(s string) manyAPI00957 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00957(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00957 fpFutureGeneric[string]

func (f manyAPI00957) next(s string) manyAPI00958 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00958(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00958 fpFutureGeneric[string]

func (f manyAPI00958) next(s string) manyAPI00959 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00959(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00959 fpFutureGeneric[string]

func (f manyAPI00959) next(s string) manyAPI00960 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00960(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00960 fpFutureGeneric[string]

func (f manyAPI00960) next(s string) manyAPI00961 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00961(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00961 fpFutureGeneric[string]

func (f manyAPI00961) next(s string) manyAPI00962 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00962(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00962 fpFutureGeneric[string]

func (f manyAPI00962) next(s string) manyAPI00963 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00963(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00963 fpFutureGeneric[string]

func (f manyAPI00963) next(s string) manyAPI00964 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00964(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00964 fpFutureGeneric[string]

func (f manyAPI00964) next(s string) manyAPI00965 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00965(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00965 fpFutureGeneric[string]

func (f manyAPI00965) next(s string) manyAPI00966 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00966(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00966 fpFutureGeneric[string]

func (f manyAPI00966) next(s string) manyAPI00967 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00967(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00967 fpFutureGeneric[string]

func (f manyAPI00967) next(s string) manyAPI00968 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00968(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00968 fpFutureGeneric[string]

func (f manyAPI00968) next(s string) manyAPI00969 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00969(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00969 fpFutureGeneric[string]

func (f manyAPI00969) next(s string) manyAPI00970 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00970(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00970 fpFutureGeneric[string]

func (f manyAPI00970) next(s string) manyAPI00971 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00971(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00971 fpFutureGeneric[string]

func (f manyAPI00971) next(s string) manyAPI00972 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00972(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00972 fpFutureGeneric[string]

func (f manyAPI00972) next(s string) manyAPI00973 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00973(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00973 fpFutureGeneric[string]

func (f manyAPI00973) next(s string) manyAPI00974 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00974(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00974 fpFutureGeneric[string]

func (f manyAPI00974) next(s string) manyAPI00975 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00975(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00975 fpFutureGeneric[string]

func (f manyAPI00975) next(s string) manyAPI00976 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00976(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00976 fpFutureGeneric[string]

func (f manyAPI00976) next(s string) manyAPI00977 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00977(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00977 fpFutureGeneric[string]

func (f manyAPI00977) next(s string) manyAPI00978 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00978(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00978 fpFutureGeneric[string]

func (f manyAPI00978) next(s string) manyAPI00979 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00979(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00979 fpFutureGeneric[string]

func (f manyAPI00979) next(s string) manyAPI00980 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00980(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00980 fpFutureGeneric[string]

func (f manyAPI00980) next(s string) manyAPI00981 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00981(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00981 fpFutureGeneric[string]

func (f manyAPI00981) next(s string) manyAPI00982 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00982(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00982 fpFutureGeneric[string]

func (f manyAPI00982) next(s string) manyAPI00983 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00983(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00983 fpFutureGeneric[string]

func (f manyAPI00983) next(s string) manyAPI00984 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00984(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00984 fpFutureGeneric[string]

func (f manyAPI00984) next(s string) manyAPI00985 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00985(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00985 fpFutureGeneric[string]

func (f manyAPI00985) next(s string) manyAPI00986 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00986(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00986 fpFutureGeneric[string]

func (f manyAPI00986) next(s string) manyAPI00987 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00987(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00987 fpFutureGeneric[string]

func (f manyAPI00987) next(s string) manyAPI00988 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00988(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00988 fpFutureGeneric[string]

func (f manyAPI00988) next(s string) manyAPI00989 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00989(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00989 fpFutureGeneric[string]

func (f manyAPI00989) next(s string) manyAPI00990 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00990(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00990 fpFutureGeneric[string]

func (f manyAPI00990) next(s string) manyAPI00991 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00991(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00991 fpFutureGeneric[string]

func (f manyAPI00991) next(s string) manyAPI00992 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00992(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00992 fpFutureGeneric[string]

func (f manyAPI00992) next(s string) manyAPI00993 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00993(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00993 fpFutureGeneric[string]

func (f manyAPI00993) next(s string) manyAPI00994 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00994(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00994 fpFutureGeneric[string]

func (f manyAPI00994) next(s string) manyAPI00995 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00995(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00995 fpFutureGeneric[string]

func (f manyAPI00995) next(s string) manyAPI00996 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00996(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00996 fpFutureGeneric[string]

func (f manyAPI00996) next(s string) manyAPI00997 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00997(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00997 fpFutureGeneric[string]

func (f manyAPI00997) next(s string) manyAPI00998 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00998(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00998 fpFutureGeneric[string]

func (f manyAPI00998) next(s string) manyAPI00999 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00999(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

type manyAPI00999 fpFutureGeneric[string]

func (f manyAPI00999) next(s string) manyAPI00000 {
	pb := func(*msgBuilder) error {
		if s == "666" {
			return errDummy
		}
		return nil
	}
	 return manyAPI00000(callFpFutureGeneric[string, string](fpFutureGeneric[string](f), uint64(100), uint16(1000), pb))
}

func BenchmarkManyTypesGeneric(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		last := manyAPI00000{pipe: &pipeline{steps: make([]pipelineStep, 0, b.N)}}.
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
			next("fooo").
		next("fooo")
		if last.stepIndex != 999 { b.Fatalf("wrong: %d", last.stepIndex) }
	}
}
