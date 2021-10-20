package predefined

const (
	AuditStatusInProgress = iota //进程中
	AuditStatusReject            //拒绝
	AuditStatusPass              //通过
)

const (
	ProfileInfoSexFemale = iota //女
	ProfileInfoSexMale          //男
)

const (
	ProfileInfoBasicMarriageUnmarried = iota //未结婚
	ProfileInfoBasicMarriageDivorced         //离异
	ProfileInfoBasicMarriageMarried          //已婚
)

const (
	ProfileInfoAddressTypeHometown  = iota //家乡
	ProfileInfoAddressTypeResidence        //居住地
)

const (
	ProfileInfoConstellationAries       = iota //白羊座：3月21日～4月20日 (Aries)
	ProfileInfoConstellationTaurus             //金牛座：4月21日～5月21日 (Taurus)
	ProfileInfoConstellationGemini             //双子座：5月22日～6月21日 (Gemini)
	ProfileInfoConstellationCancer             //巨蟹座：6月22日～7月22日 (Cancer)
	ProfileInfoConstellationLeo                //狮子座：7月23日～8月23日 (Leo)
	ProfileInfoConstellationVirgo              //处女座：8月24日～9月23日 (Virgo)
	ProfileInfoConstellationLibra              //天秤座：9月24日～10月23日 (Libra)
	ProfileInfoConstellationScorpio            //天蝎座：10月24日～11月22日 (Scorpio)
	ProfileInfoConstellationSagittarius        //射手座：11月23日～12月21日 (Sagittarius)
	ProfileInfoConstellationCapricorn          //摩羯座：12月22日～1月20日 (Capricorn)
	ProfileInfoConstellationAquarius           //水瓶座：1月21日～2月19日 (Aquarius)
	ProfileInfoConstellationPisces             //双鱼座：2月20日～3月20日 (Pisces)
)

const (
	ProfileInfoProfessionAnnualIncome5_15    = iota //年薪：5万~15万
	ProfileInfoProfessionAnnualIncome15_30          //年薪：15万~30万
	ProfileInfoProfessionAnnualIncome30_50          //年薪：30万~50万
	ProfileInfoProfessionAnnualIncome50_100         //年薪：50万~100万
	ProfileInfoProfessionAnnualIncome100_500        //年薪：100万~500万
	ProfileInfoProfessionAnnualIncome500_           //年薪：500万以上
)

const (
	ProfileQuestionTypeValues = iota //价值观
)

const (
	ProfileInfoIntroductionTypeAboutMe          = iota //关于我
	ProfileInfoIntroductionTypeFamilyBackground        //家庭背景
	ProfileInfoIntroductionTypeHobbies                 //兴趣爱好
)

const (
	ProfileProofProfessionTypeCompanySocialSecurity           = iota //社保
	ProfileProofProfessionTypeCompanyEnterpriseOfficeSoftware        //企业软件
	ProfileProofProfessionTypeCompanyLicense                         //执照
	ProfileProofProfessionTypeCompanyWorkPermit                      //工作证
	ProfileProofProfessionTypeCompanyPaySlip                         //工资单
	ProfileProofProfessionTypeCompanyOffer                           //录取Offer
	ProfileProofProfessionTypeStudent                                //学生
)

const (
	ProfileProofEducationTypeCHSI         = iota //学信网
	ProfileProofEducationTypeDiplomaImage        //毕业证文凭照片
	ProfileProofEducationTypeDiplomaID           //毕业证文凭编号
	ProfileProofEducationTypeCSCSE               //教育部留学服务中心认证
	ProfileProofEducationTypeOldCSCSE            //教育部留学服务中心认证（旧版）
)
