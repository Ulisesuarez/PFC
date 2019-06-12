<html>
    <head>
        <meta charset="utf-8">
        <script type="text/javascript" src="js/newSearch.js"></script>
        <script type="text/javascript" src="js/LoadLanguage.js"></script>
        <script type="text/javascript" src="js/lib/jquery-gentleSelect.js"></script>
        <link rel="stylesheet" type="text/css" href="css/jquery.alerts.css" />
        <link rel="stylesheet" href="css/font-awesome/css/font-awesome.min.css">
        <script type="text/javascript" src="js/lib/jquery.alerts.js"></script>
        <!-- <script type="text/javascript" src="js/lib/jquery.multiselect.js"></script>
        <link href="css/jquery.multiselect.css" rel="stylesheet" type="text/css"> -->
        <script type="text/javascript">
            $(document).ready(function() {
                $('button').button();
                $('.bdate').datepicker();
                $('.fsdate').datepicker();
                $('.tsdate').datepicker();
                
                $('#<%=request.getParameter("tabName")%>_modalities').gentleSelect({
			   		 columns: 3,
		    		 itemWidth: 20,
				     maxDisplay: 3,				    
				     prompt: "ALL", 
		   		 });	
                //     $('#<%=request.getParameter("tabName")%>_modalities').multiselect({
    	        // 	columns: 2,
    	        // 	minWidth: 157,
    	        // 	selectedList: 4,
    	        // 	header: true,
    	        // 	includeSelectAllOption: true,
    	        //     noneSelectedText: "ALL"
    	        // });
//                var lang = getCookie('language');
//                if (lang == 'none') {
//                $.getScript('js/i18n/Bundle.js', function() {
//                    loadLabels();
//                });
//                } else {
//                var fileName = 'js/i18n/' + "Bundle_" + lang + ".js";
//                        $.getScript(fileName, function() {
//                    loadLabels();
//                });
//                }
//                $("label[for=patId]").text(languages['PatientId']);
//                $("label[for=patName]").text(languages['PatientName']);
//                $("label[for=accessionNo]").text(languages['AccessionNumber']);
//                $("label[for=birthDate]").text(languages['BirthDate']);
//                $("label[for=studyDesc]").text(languages['StudyDesc']);
//                $("label[for=referPhysician]").text(languages['ReferPhysician']);
//                $(".bdate").prev().text(languages['BirthDate']);
//                $(".fsdate").prev().text(languages['FromStudyDate']);
//                $(".tsdate").prev().text(languages['ToStudyDate']);
//                $(".searchBtn").text(languages['Search']);
//                $(".clearBtn").text(languages['Reset']);

            });  //for document.ready
        </script>

        <style>
            /* css for timepicker */
            .ui-timepicker-div .ui-widget-header { margin-bottom: 8px; }
            .ui-timepicker-div dl { text-align: left; }
            .ui-timepicker-div dl dt { height: 25px; margin-bottom: -25px; }
            .ui-timepicker-div dl dd { margin: 0 10px 10px 65px; }
            .ui-tpicker-grid-label { background: none; border: none; margin: 0; padding: 0; }

            #newQueryParamDiv label {
                float: left;
                width: 10em;
                font-size: 13px;
            }            
            
/*             .ui-state-default { */
/*             	border: 2px solid #2A2A2A !important; */
/*             } */
            
            .gentleselect-label {
            	float: left;
    			/*font-weight: bold;*/
    
    			border: 2px solid #2A2A2A !important;
    			background-image: url('css/images/expand.png');
    			display: block;    
				width: 13em;
				height: 20px;
				text-align: center;
				line-height: 2em;
				padding: 0px !important;
            }
            
           /* .gentleselect-dialog > ul {
            	font-weight: bold;
            }*/
            
            .gentleselect-dialog > ul > li.selected {
            	color: #fff;
			    background-color: #0b93d5;
            }
            
            .gentleselect-dialog > ul > li:hover {
            	 background-color: #0972a5;
			    color: #fff;
            }
            #findTable label{
                font-size: 14px;
            }
           
        </style>

    </head>
    <body>
        <div id="newQueryParamDiv" class="ui-widget-content">
<!--             <table id="findTable" class="ui-widget-content" width="930px" height="90%" cell-spacing="4" style="margin-left: 10px; font-size: 100%;border: none;"> -->
<table id="findTable" class="ui-widget-content" width="100%" height="90%"  cellspacing="4" style="margin-left: 10px; font-size: 100%;border: none;">
                <tr>
<!--                     <td style="width:50px;"> -->
						<td>
                        <label for="patId">Patient ID</label>
                        <input type="search" id="patId" placeholder="Patient ID" />
                    </td>
                    <td>
                        <label for="patName">Patient Name</label>
                        <input type="search" id="patName" placeholder="Patient Name" />
                    </td>
                   <!-- <td>
                        <label for="accessionNo">Accession Number</label>
                        <input type="search" id="accessionNo" />
                    </td>
                    <td>
                        <label>Birth Date</label>
                        <input type="search" class="bdate" />
                    </td>-->
                    
                    <td>
                        <label>Study Date (From)</label>
                        <input type="search" class="fsdate" placeholder="MM/DD/YYYY" />
                    </td>
                    <td>
                        <label>Study Date (To)</label>
                        <input type="search" class="tsdate" placeholder="MM/DD/YYYY"/>
                    </td>
                </tr>
                <tr>   
                	<td>
                        <label for="studyDesc">Study Description</label>
                        <input type="search" id="studyDesc" placeholder="Study/Procedure name" />
                    </td>
                                     
                    <td>
                        <label for="referPhysician">Referring Physician</label>
                        <input type="search" id="referPhysician" placeholder="Referring Physician" />
                    </td>
                    <td>
                    	<div style="float: left; width: 10em; font-size: 14px; margin-top: 6px;"> Modality </div>
                       
                        <select id="<%=request.getParameter("tabName")%>_modalities" multiple="multiple" class="modalitiesList <%=request.getParameter("tabName")%>_modalities" style="display: none;">
                            
                            <option value="CT">CT</option>
                            <option value="CR">CR</option>
                            <option value="XA">XA</option>
                            <option value="MR">MR</option>
                            <option value="US">US</option>
                            <option value="DX">DX</option>
                            <option value="PX">PX</option>
                            <option value="SC">SC</option>
                            <option value="NM">NM</option>
                            <option value="RF">RF</option>
                            <option value="OT">OT</option>
                            <option value="ECG">ECG</option>  
                                                                             
                        </select>    

                    </td>                     
                    <td colspan="2">   
                    	<button id="okBtn" onclick="searchClick(this)" style="float: left;width:7em;height:27px; font-size: 14px; margin-right: 10px;" class="searchBtn">Search</button>                                            	
    	                	<button id="resetBtn" onclick="resetClick(this, <%=request.getParameter("tabName")%>_modalities);" style="float: left;width:7em; height:27px; font-size: 14px;" class="clearBtn"> Reset</button>           	                    	
	                	</td>
                </tr>
            </table>
        </div>
    </body>
</html>
