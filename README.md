# Enigma Cryptanalysis

## Answer for given [ct.txt](https://github.com/shreyas-sriram/enigma-cryptanalysis/blob/main/ct.txt)

```
Gamma VI IV III 
D A B Q 
AG BN DI ER FY HJ KU LW MS PQ
```

## Creating test cases

1. Generate encrypted text using [github.com/emedvedev/enigma](https://github.com/emedvedev/enigma)
```
enigma "EMMAWOODHOUSEXHANDSOMEXCLEVERXANDRICHXWITHACOMFORTABLEHOMEANDHAPPYDISPOSITIONXSEEMEDTOUNITESOMEOFTHEBESTBLESSINGSOFEXISTENCEXANDHADLIVEDNEARLYTWENTYXONEYEARSINTHEWORLDWITHVERYLITTLETODISTRESSORVEXHERXSHEWASTHEYOUNGESTOFTHETWODAUGHTERSOFAMOSTAFFECTIONATEXINDULGENTFATHERXANDHADXINCONSEQUENCEOFHERSISTERXSMARRIAGEXBEENMISTRESSOFHISHOUSEFROMAVERYEARLYPERIODXHERMOTHERHADDIEDTOOLONGAGOFORHERTOHAVEMORETHANANINDISTINCTREMEMBRANCEOFHERCARESSESXANDHERPLACEHADBEENSUPPLIEDBYANEXCELLENTWOMANASGOVERNESSXWHOHADFALLENLITTLESHORTOFAMOTHERINAFFECTIONXSIXTEENYEARSHADMISSTAYLORBEENINMRXWOODHOUSEXSFAMILYXLESSASAGOVERNESSTHANAFRIENDXVERYFONDOFBOTHDAUGHTERSXBUTPARTICULARLYOFEMMAXBETWEENTHEMITWASMORETHEINTIMACYOFSISTERSXEVENBEFOREMISSTAYLORHADCEASEDTOHOLDTHENOMINALOFFICEOFGOVERNESSXTHEMILDNESSOFHERTEMPERHADHARDLYALLOWEDHERTOIMPOSEANYRESTRAINTXANDTHESHADOWOFAUTHORITYBEINGNOWLONGPASSEDAWAYXTHEYHADBEENLIVINGTOGETHERASFRIENDANDFRIENDVERYMUTUALLYATTACHEDXANDEMMADOINGJUSTWHATSHELIKEDXHIGHLYESTEEMINGMISSTAYLORXSJUDGMENTXBUTDIRECTEDCHIEFLYBYHEROWNXTHEREALEVILSINDEEDOFEMMAXSSITUATIONWERETHEPOWEROFHAVINGRATHERTOOMUCHHEROWNWAYXANDADISPOSITIONTOTHINKALITTLETOOWELLOFHERSELFXTHESEWERETHEDISADVANTAGESWHICHTHREATENEDALLOYTOHERMANYENJOYMENTSXTHEDANGERXHOWEVERXWASATPRESENTSOUNPERCEIVEDXTHATTHEYDIDNOTBYANYMEANSRANKASMISFORTUNESWITHHERXSORROWCAMEXXAGENTLESORROWXXBUTNOTATALLINTHESHAPEOFANYDISAGREEABLECONSCIOUSNESSXMISSTAYLORMARRIEDXITWASMISSTAYLORXSLOSS" --rotors "Gamma I IV III" --rings="1 1 1 16" --position "F U B Q" --plugboard "AQ BN DI ER FY HJ KU LW MS PG" --reflector "C-thin"
```

2. Store encrypted text in a file
```
echo -n "QYQXZXEHIGFISNYLHCDBRQAEIQQJLRSCAGABUZIDBUEKIKJZDCTWWQBMEPWDTJXBEUVXXRRGTAALQNUOBRRZKYBBYRAJBOHWHNPTWQNXGSWHTSUHAHDBEXFRMDNVDNUFFCOXCRBQFJGUHSGTQWOJFUUNWXUOODRJWDIHJBWOZXNJWJRZDQEHLBEVDHLKUYNDIRSRNWKFWATAXHZDOFNBOKUEWAHVBAQQXHNPYDEFHRELCVYWWBUDIMDFEVBIZFSKXERZJSSMOSDWUOPOKOYAJSEOIWMXXVUGLZYDVYKVOWKTOZKZESZZMORDYTCZSNKYVTUBAGRJVTXJEDGCLVMGSBQHODCKJIBLCENJJULKSQWQVFGIOWMWAJBBMRQOWPHALSTWIYMASNTCFAARFLOBRCLCNTRAJCCRYAPIEGKGIJFCXOOVFQBNTPPWELBOKFNDVCYRNOYQODQSIMCOLAYZLYTHGOSWGZYDKGLZZSUCVWVHORYDAJXUJQOJSHJNWBCPXAKUYBSNBTWSXLYVJRXAPHQLBRTCAMTDGQOWQIGMSXTPZIBQPIYEGKDVTJGACJWYWLBPVCDPKJECMSFOAVIEBJBCRIQKTDYUULUBXVADZYRKQZOFREDOVSRQBCRKYOTZVUWYNZRSIHRAVXKZQACUGDWSXRFFYDJSBAPVLFONAJYPUTIGLLOTPDFNHLHJULUVCRSDLWJOWJLSYBKXSZNHHDNJJTYAPALNSZSHAYKEFSSWKJWZKWFIBYIYIANPVRZNGIDQRKGPBWBAYXYTYNLYVUWCWNFNOVWYPDUCBJZIGGHMSESGINXOUUTDHDSGSHZYKGLGEWVIGSJDKVGOIJFIPEDCWTWTFIGTFETIDQVMKKJHLRBGBIULAVRZJXFPXCGYQFJGPLFEXHCWJPAOXOBMYEWMPMKVZKNVAZFHUVYUHVEKJZWOTXFDVDXSXTALMKPBNIAPRFGBMWSAHBONJQFFZMIJAUXFTOJEVOYGVFCBDERJUGHQPZBVQQTKBSRTUVQGOZOIFEJRSSFOGRVTHZUVPLMAZDTRIPGSZVVZWMMFKLFPUWAVZDFPQNHHMPUGQRVUODLBUZPOPYDLYEYCFOQZCUFMEBTEULIWWYTTNXDXBIXPOCKRWFKBGBUPJMYKTDXZBATAOWRRINTQIATASFJCHUMRCEWAKGBOYSXNAVHLEASNGRPAYYZAFXKDCMYGGDNOUKPRBCABGCQTIVSHNUPILNKCXRMAGUNHSVIRTKTJGUKOSJLZPYIOCXUSIUVOFETVZRWBZKFKJBKMWGCDOPXCONEXYJQKUVLPIOMQWGGELYLAZEWCEJRLYVNVYATQOJICBSCDYSNDBOWCNBQNUUGJYSOEHCMAJLUFUMOFATRWCATOIMLPHXMPPPWDTBEDMKSIKFQFGNBCAOXLKXNQOVREDZO" > a.txt
```

3. Compile and run `hillclimb`
```
go build *.go && time ./hillclimb a.txt
```
