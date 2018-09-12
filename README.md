# Image to ASCII
[![GoDoc][1]][2]
[![Build Status][3]][4]
[![Go Report Card][5]][6]

[1]: https://godoc.org/github.com/zyxar/image2ascii?status.svg
[2]: https://godoc.org/github.com/zyxar/image2ascii
[3]: https://travis-ci.org/zyxar/image2ascii.svg?branch=master
[4]: https://travis-ci.org/zyxar/image2ascii
[5]: https://goreportcard.com/badge/github.com/zyxar/image2ascii
[6]: https://goreportcard.com/report/github.com/zyxar/image2ascii

```
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMX0kxdolll,;::ccccllox0XWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMXOoc:cclllood:cdddooooolccccldk0XMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNxc:cloddddoooodc;doooooooodddoolccccdXMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMXd;coddoooooooooool;doooooooooooooddddoc:l0MMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMWkl:codooooooooooooool;dooooooooooooooooooddl::kNMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMW0o;clddooooooooooooododl;dodddddddoooooooooooodo,::lOMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMWo:c;:doooodddddooooolllllc,lccllllllllooooodddodl,odol;cMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMl;ddl;dddollcc:::::::::::::::cccc::cccc::::ccccll;odooox,xMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMWXKKk.dood;::;:cccccccllodddddddddddddxdlc::::;;;;:c:;;:lddool.kO0XMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMM0lodkkk..do:;:loc;;clooooc:;:ldooooooodo;;cxOKXXK0xc'.:lcc:,cdd.:K0OdcOMMMMMMMMMMMMMM
MMMMMMMMMMMMMMx,KMMMMl,d:,coo:.,xXWMMMMMMWKx;,odooood:'xNMMMMMMMMMMMO:ck:lo;,oo.x0XMN,dMMMMMMMMMMMMM
MMMMMMMMMMMMMN.WMW'  ,d,;do:c.kMMMMMMMMMMMMMMk.ldood;,NMMMMMMMMMMMMMMMllW::x:'oo   oMW.WMMMMMMMMMMMM
MMMMMMMMMMMMMX.MMWc .x,;xl;K'KMMMMMMMMMMMKxdkMK.oool.WMMMMMMMMMMNc...oM'kM,cx:'xc oNM0.WMMMMMMMMMMMM
MMMMMMMMMMMMMMo:KMM.ol.do'Nk;MWMMMMMMMMM;    .X;,dd;,MMMMMMMMMMM.   ; ox:Mx,dd'co.X0o:XMMMMMMMMMMMMM
MMMMMMMMMMMMMMMK:,;.c.:dc;Md;MWMMMMMMMMW.  .l 0:,dol.WMMMMMMMMMMx   ,.X,kMO'dd:.:',.,MMMMMMMMMMMMMMM
MMMMMMMMMMMMMMWMo,ooo.cdc;MN.0MMMMMMMMMMNl'';K0.dddx;,NMMMMMMMMMMN0O0Ml:MMl;ddc,xdxl;MWMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMX;cld'cdd,OMX,dWMMMMMMMMMMMMMk.lc:;;c;,kWMMMMMMMMMMM0:dMMd,dod;'lcc'kMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMK.::,'ddo,dWMxlxKNMMMMMMW0d;,,       ..,lxOKXXXKOdodXNxc:dodc'ccc:,MMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMo;dod:'ldd::o0Kkdooool:,',''ldo;...,cOkc.;;,'.',:lddl:coddl;,odddc:MMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMo;doodl,;lddlccloolc::::c,,XMMMMMMMMMMMMO.lllooooooodddlc:l.oooooo;oMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMM0'doooodl,,ccllllllcccclkd,KNX0OkkkkOKX0clXkdlcccccccccokXM:cdooood.0MMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMW'ooooood,0NOkddddxkOKNMMMKx'ckOO.O0Oc.k0MMMMMWNXXXXXNMMMMM0'oooood,xMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMo;dooodc;MMMMMMMMMMMMMMMMMM'KMMM'MMMX'MMMMMMMMMMMMMMMMMMMMMd;doood'KMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMo;dooodc:MMMMMMMMMMMMMMMMMM.0MM0.KWMx'MMMMMMMMMMMMMMMMMMMMMO'dood;lMMMMMMMMMMMMMMM
MNkWMMMMMMMMMMMMMN'ooooodo.NMMMMMMMMMMMMMMMMMOldookdddlXMMMMMMMMMMMMMMMMMMMMM::dod::WMMMMMMMMMMMMWWM
dc;;WMMMMMMMMMMMM;cdooddl,kMMMMMMMMMMMMMMMMMMMMWMMMMMMMMMMMMMMMMMMMMMMMMMMMM0'dool'MMMMMMMMMMMKxxxko
c,k,kMMMMMMMMMWNl:doooc'lXMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMO'loddcckXMMMMMMKl,.ldxx
Mx,c;OOkdlccccc:;ccccck.0MWMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMKo::coocccxNMWx;;cx0NWM
MMKoc;;;:oooooodk0XNNMM;xMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMd'kc:cll::lcc;lWMMMMM
MMMMWNNWMMMMMMMMMMMMMMM;oMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMx;MMXOdool:ccoMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMO0MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMk:MMMMMMMWNXWMMMMMMMM
```


inspired by [jp2a](http://csl.sublevel3.org/jp2a/)


## Install

`go get github.com/zyxar/image2ascii`


## Supported [Formats](https://en.wikipedia.org/wiki/Image_file_formats)

- [x] `jpeg`
- [x] `png`
- [x] `gif`
- [x] `bmp`
- [x] `tiff`
- [x] `ico`
- [x] `vp8l`
- [x] `webp`


## Cmd

```
Usage: image2ascii [options] {IMAGE FILE}
  -c  enable colour mode
  -h int
      set image height
  -i  enable invert mode
  -w int
      set image width
  -x  enable flip-x mode
  -y  enable flip-y mode
```


## Alternative

[node-jp2a](https://github.com/zyxar/node-jp2a): `npm install -g node-jp2a`


## License
[Apache 2.0](http://opensource.org/licenses/Apache-2.0)
