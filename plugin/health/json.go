package health

type healthId struct {
	Rows []struct {
		HealthID                 string `json:"healthId"`
		SubmitTime               string `json:"submitTime,omitempty"`
		DataStatus               string `json:"dataStatus"`
		DataType                 string `json:"dataType"`
		DataDate                 string `json:"dataDate"`
		CounsellorApprovalStatus string `json:"counsellorApprovalStatus,omitempty"`
		CounsellorApprovalView   string `json:"counsellorApprovalView,omitempty"`
		SecretaryApprovalStatus  string `json:"secretaryApprovalStatus,omitempty"`
		SecretaryApprovalView    string `json:"secretaryApprovalView,omitempty"`
		EndTime                  string `json:"endTime"`
	} `json:"rows"`
	Total int `json:"total"`
}

type rosterId struct {
	IsSuccess bool `json:"isSuccess"`
	Data      struct {
		Model string `json:"model"`
	} `json:"data"`
}

type rosterIdModel struct {
	Rows []struct {
		IsNormal          int    `json:"isNormal"`
		StudentCode       string `json:"studentCode"`
		SchoolDeptID      string `json:"schoolDeptId"`
		HadInit           int    `json:"hadInit"`
		IsHomeHubei       int    `json:"isHomeHubei"`
		State             int    `json:"state"`
		HomeAddress       string `json:"homeAddress"`
		SchoolDeptName    string `json:"schoolDeptName"`
		MajorID           string `json:"majorId"`
		PoliticalStatusID string `json:"politicalStatusId"`
		NationID          string `json:"nationId"`
		SchoolAreaName    string `json:"schoolAreaName"`
		TimeIntervalMor   int    `json:"timeIntervalMor"`
		Grade             int    `json:"grade"`
		StudentName       string `json:"studentName"`
		Sex1              string `json:"sex1"`
		StudentType       string `json:"studentType"`
		Dormitory         string `json:"dormitory"`
		StudentID         string `json:"studentId"`
		ExecutiveClassID  string `json:"executiveClassId"`
		DormitoryNum      string `json:"dormitoryNum"`
		TeacherName       string `json:"teacherName"`
		HomeAddressJSON   struct {
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
			Street   string `json:"street"`
		} `json:"homeAddressJson"`
		Sex                string `json:"sex"`
		NativePlaceName    string `json:"nativePlaceName"`
		IsPostGraduate     int    `json:"isPostGraduate"`
		RosterID           string `json:"rosterId"`
		IsReturnSchool     int    `json:"isReturnSchool"`
		CreateTime         string `json:"createTime"`
		ExecutiveClassName string `json:"executiveClassName"`
		LastModifyTime     string `json:"lastModifyTime"`
		MajorName          string `json:"majorName"`
		SourceName         string `json:"sourceName"`
	} `json:"rows"`
	Total int `json:"total"`
}

type ChickID struct {
	Rows []struct {
		ChickID           string `json:"chickId"`
		Address           string `json:"address"`
		CheckTime         string `json:"checkTime"`
		CheckingEndTime   string `json:"checkingEndTime"`
		CheckingBeginTime string `json:"checkingBeginTime"`
		SysMonitoryStatus int    `json:"sysMonitoryStatus"`
		Remark            string `json:"remark"`
		Type              string `json:"type"`
		SysMonitoryName   string `json:"sysMonitoryName"`
		CheckDate         string `json:"checkDate"`
		Status            string `json:"status"`
	} `json:"rows"`
	Total int `json:"total"`
}
