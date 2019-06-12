<!--
/* ***** BEGIN LICENSE BLOCK *****
* Version: MPL 1.1/GPL 2.0/LGPL 2.1
*
* The contents of this file are subject to the Mozilla Public License Version
* 1.1 (the "License"); you may not use this file except in compliance with
* the License. You may obtain a copy of the License at
* http://www.mozilla.org/MPL/
*
* Software distributed under the License is distributed on an "AS IS" basis,
* WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
* for the specific language governing rights and limitations under the
* License.
*
* The Original Code is part of Oviyam, an web viewer for DICOM(TM) images
* hosted at http://skshospital.net/pacs/webviewer/oviyam_0.6-src.zip
*
* The Initial Developer of the Original Code is
* Raster Images
* Portions created by the Initial Developer are Copyright (C) 2014
* the Initial Developer. All Rights Reserved.
*
* Contributor(s):
* Babu Hussain A
* Devishree V
* Meer Asgar Hussain B
* Prakash J
* Suresh V
* Yogapraveen K
*
* Alternatively, the contents of this file may be used under the terms of
* either the GNU General Public License Version 2 or later (the "GPL"), or
* the GNU Lesser General Public License Version 2.1 or later (the "LGPL"),
* in which case the provisions of the GPL or the LGPL are applicable instead
* of those above. If you wish to allow use of your version of this file only
* under the terms of either the GPL or the LGPL, and not to allow others to
* use your version of this file under the terms of the MPL, indicate your
* decision by deleting the provisions above and replace them with the notice
* and other provisions required by the GPL or the LGPL. If you do not delete
* the provisions above, a recipient may use your version of this file under
* the terms of any one of the MPL, the GPL or the LGPL.
*
* ***** END LICENSE BLOCK ***** */
-->

<%@page language="java" contentType="text/html; charset=UTF-8"
pageEncoding="UTF-8" isELIgnored="false"%>
    <%@page errorPage="ErrorPage.jsp"%>
        <%@taglib prefix="ser" uri="/WEB-INF/tlds/SeriesDetails.tld"%>
            <%@taglib prefix="img" uri="/WEB-INF/tlds/ImageInfo.tld"%>
                <%@taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core"%>
                    <%@taglib prefix="fmt" uri="http://java.sun.com/jsp/jstl/fmt"%>
                        <%@taglib uri="http://java.sun.com/jsp/jstl/functions" prefix="fn"%>

                            <%
String patName = new String(request.getParameter("patientName")
        .getBytes("ISO-8859-1"), "UTF-8");
String studyDesc = new String(request.getParameter("studyDesc")
        .getBytes("ISO-8859-1"), "UTF-8");
%>

                                <html>

                                <head>

                                    <style type="text/css">
                                        .heading1 {
                                            font-family: Arial;
                                            font-size: 15px;
                                            font-weight: bold;
                                            padding: 3px;
                                            width: 100%;
                                            white-space: nowrap;
                                            text-overflow: ellipsis;
                                            overflow: hidden;
                                        }
                                        
                                        .seriesTable {
                                            table-layout: fixed;
                                            width: 100%;
                                            font-family: Arial;
                                            font-size: 12px;
                                            border-spacing: 0px;
                                            padding-left: 2px;
                                        }
                                        
                                        .scale-image {
                                            width: 30%;
                                            height: 30%;
                                        }
                                        
                                        .accordion1 {
                                            padding: 5px;
                                            font-weight: bold;
                                            background-color: #666;
                                            background-position: right center;
                                            background-repeat: no-repeat;
                                            color: #F90;
                                            border: 1px solid #333;
                                            font-size: 14px;
                                            font-family: Arial;
                                        }
                                        
                                        .accordion1:hover {
                                            background-color: #555;
                                        }
                                        
                                        #otherStudiesInfo {
                                            padding: 5px;
                                            font-weight: bold;
                                            cursor: pointer;
                                            background-color: #222;
                                            color: #fff;
                                            font-family: Arial;
                                            font-size: 14px;
                                            text-align: center;
                                            margin: 7px 0px 7px 0;
                                        }
                                        
                                        .open1 {
                                            background-image: url(images/arrow_up.gif);
                                        }
                                        
                                        .close1 {
                                            background-image: url(images/arrow_down.gif);
                                            width: 100%;
                                        }
                                        
                                        #previews::-webkit-scrollbar {
                                            width: 14px;
                                            background: #464646;
                                        }
                                        
                                        #previews::-webkit-scrollbar-track {
                                            -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
                                        }
                                        
                                        #previews::-webkit-scrollbar-thumb {
                                            -webkit-border-radius: 10px;
                                            border-radius: 10px;
                                            background: #262626;
                                            -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.5);
                                        }
                                        
                                        .button {
                                            background-color: #FF8A00;
                                            border: none;
                                            color: white;
                                            padding: 3px 8px;
                                            text-align: center;
                                            text-decoration: none;
                                            display: inline-block;
                                            font-size: 12px;
                                            margin-left: 10px;
                                            margin-right: 10px;
                                            margin: 2px 1px;
                                            cursor: pointer;
                                            -webkit-transition-duration: 0.4s;
                                            /* Safari */
                                            transition-duration: 0.4s;
                                        }
                                        
                                        .button1 {
                                            background-color: #FF8A00;
                                            /* Green */
                                            border: none;
                                            color: white;
                                            padding: 3px 20px;
                                            text-align: center;
                                            text-decoration: none;
                                            display: inline-block;
                                            font-size: 12px;
                                            margin-left: 10px;
                                            margin-right: 10px;
                                            margin: 15px 0px;
                                            cursor: pointer;
                                            -webkit-transition-duration: 0.4s;
                                            /* Safari */
                                            transition-duration: 0.4s;
                                        }
                                        
                                        .preview {
                                            overflow: auto;
                                            border: 1px solid #FF8A00;
                                            height: 90%;
                                            padding: 5px;
                                            background: #333
                                        }
                                        
                                        .modal-confirm {
                                            display: none;
                                            position: fixed;
                                            z-index: 1;
                                            left: 0;
                                            top: 0;
                                            width: 100%;
                                            height: 100%;
                                            overflow: auto;
                                            background-color: #666;
                                            background-color: rgba(0, 0, 0, 0.4);
                                        }
                                        /* Modal Content */
                                        
                                        .modal-content-confirm {
                                            background-color: #666;
                                            margin: auto;
                                            border: 1px solid chocolate;
                                            margin-top: 25%;
                                            width: 575px;
                                            height: 140px;
                                        }
                                        
                                        .loader {
                                            border: 4px solid #f3f3f3;
                                            border-radius: 50%;
                                            border-top: 4px solid #FF8A00;
                                            border-bottom: 4px solid #FF8A00;
                                            width: 40px;
                                            height: 40px;
                                            -webkit-animation: spin 2s linear infinite;
                                            animation: spin 2s linear infinite;
                                        }
                                        
                                        @-webkit-keyframes spin {
                                            0% {
                                                -webkit-transform: rotate(0deg);
                                            }
                                            100% {
                                                -webkit-transform: rotate(360deg);
                                            }
                                        }
                                        
                                        @keyframes spin {
                                            0% {
                                                transform: rotate(0deg);
                                            }
                                            100% {
                                                transform: rotate(360deg);
                                            }
                                        }
                                    </style>

                                    <script type="text/javascript">
                                        var pat = $.cookies.get('patient');
                                        var confirmDownload = document.getElementById('modalConfirm');
                                        var loading = document.getElementById('loadingDiv');
                                        loading.style.display = "none";
                                        var serverURL = pat.serverURL.trim();

                                        var totalJpeg = 0;
                                        var downloadRequestJson = {};
                                        var selectedSeries = [];
                                        var patientInfo = {};
                                        var result;
                                        var requestedTime;
                                        var jpeg = document.getElementById('jpeg');
                                        var dicom = document.getElementById('dicom');
                                        var both = document.getElementById('both');

                                        var patientUID = "";
                                        var studyUID = "";
                                        var studyDate = "";
                                        var studyDesc = "";

                                        function showOtherStudies(patientID, selectedStudyId, dcmUrl, wadoUrl) {
                                            $.post("otherStudies.do", {
                                                "patientID": patientID,
                                                "studyUID": selectedStudyId,
                                                "dcmURL": dcmUrl
                                            }, function(data) {
                                                $.each(data, function(i, study) {
                                                    var div = document.createElement("div");
                                                    var url = "downloadStudy.jsp?patient=" + patientID + "&study=" + study["studyUID"] + "&dcmURL=" + dcmUrl + "&wadoUrl=" + wadoUrl + "&studyDesc=" + study["dateDesc"];
                                                    jQuery("#otherStudies1").show();
                                                    jQuery("#otherStudies1").load(url);
                                                });
                                            }, "json");
                                        }

                                        function resetSelection() {
                                            totalJpeg = 0;
                                            downloadRequestJson = {};
                                            patientInfo = {};
                                            selectedSeries = [];
                                            $('#No_Images').text("");
                                        }

                                        function addTotalImage(No_Images) {
                                            var No_image = parseInt(No_Images);
                                            totalJpeg = totalJpeg + No_image;
                                            $('#No_Images').text(totalJpeg + ' images')
                                        }

                                        function addSeries(imageInfo) {
                                            patientUID = imageInfo[0];
                                            studyDesc = imageInfo[1];
                                            studyDate = imageInfo[2];
                                            studyUID = imageInfo[3];
                                            item = {}
                                            item['seriesDesc'] = imageInfo[4].trim();
                                            item['seriesDate'] = imageInfo[5].trim();
                                            item['seriesID'] = imageInfo[6].trim();
                                            item['instanceNumber'] = imageInfo[7].trim();
                                            selectedSeries.push(item);
                                        }

                                        function addPatientInfo() {
                                            patientInfo = {
                                                patientID: patientUID,
                                                patientName: pat.pat_Name.trim(),
                                                studyID: studyUID,
                                                studyDesc: studyDesc,
                                                studyDate: studyDate,
                                                dcmURL: pat.dicomURL.trim(),
                                                serverURL: pat.serverURL.trim()
                                            }
                                        }

                                        function getDownloadRequest(downloadFileType) {
                                            downloadRequestJson = {
                                                selectedSeries: selectedSeries,
                                                patientInfo: patientInfo,
                                                downloadFileType: downloadFileType
                                            }
                                        }

                                        function closeDownloadModal() {
                                            var modal = document.getElementById('downloadModal');
                                            modal.style.display = "none";
                                        }

                                        function downloadBtnClick() {
                                            addPatientInfo();
                                            if (selectedSeries.length == 0) {
                                                $('#confirm-Download').text("Cannot Initiate Download");
                                                $('#fileInfo').text("Please select any study/Series to download");
                                                $('#fileSize').text("");
                                                document.getElementById('confirmDownloadBtn').style.visibility = "hidden";
                                                confirmDownload.style.display = "block";
                                                $('#cancelBtn').text("Ok");
                                            } else {
                                                $('#confirm-Download').text("CONFIRM DOWNLOAD");
                                                document.getElementById('confirmDownloadBtn').style.visibility = "visible";
                                                confirmDownload.style.display = "none";
                                                if (jpeg.checked) {
                                                    getDownloadRequest('JPEG');
                                                } else if (dicom.checked) {
                                                    getDownloadRequest('DICOM');
                                                } else {
                                                    getDownloadRequest('JPEG&DICOM');
                                                }
                                                var downloadRequest = JSON.stringify(downloadRequestJson);
                                                processDownload(downloadRequest);
                                                loading.style.display = "block";
                                            }
                                        }


                                        function processDownload(downloadRequest) {
                                            $.get("downloadStudy.do", {
                                                'downloadRequest': downloadRequest
                                            }, function(data) {
                                                confirmDownload.style.display = "block";
                                                loading.style.display = "none";
                                                if (data.trim() == 'ERROR') {
                                                    $('#confirm-Download').text("Error while processing download");
                                                    $('#fileInfo').text("Please try again..");
                                                    $('#fileSize').text("");
                                                    $('#cancelBtn').text("OK");
                                                    document.getElementById('confirmDownloadBtn').style.visibility = "hidden";
                                                } else {
                                                    $('#cancelBtn').text("Cancel");
                                                    $('confirmDownloadBtn').text("Download");
                                                    result = data;
                                                    var fileInfo = result.split('|');
                                                    $('#fileInfo').text(fileInfo[0]);
                                                    $('#fileSize').text(fileInfo[1]);
                                                    requestedTime = fileInfo[2];
                                                }
                                            }, 'text');
                                        }

                                        $('#confirmDownloadBtn').click(function() {
                                            var download = "yes";
                                            var outString = "downloadStudy.do?requestedTime=" + requestedTime + "&download=" + download;
                                            location.href = outString
                                        });

                                        function cancelDownload() {
                                            if (result == "") {
                                                confirmDownload.style.display = "none";
                                            } else {
                                                confirmDownload.style.display = "none";
                                                /* Calling Download.jsp using Ajax */
                                                $.ajax({
                                                    method: "POST",
                                                    url: "downloadStudy.do",
                                                    data: {
                                                        "requestedTime": requestedTime,
                                                        "download": "no"
                                                    },
                                                    async: false
                                                });
                                            }
                                        }
                                    </script>
                                </head>

                                <body>
                                    <div id="loadingDiv" class="modal-confirm">
                                        <div id="loading" style="margin-left: 50%;margin-top: 25%">
                                            <p class="loader"></p>
                                        </div>
                                    </div>

                                    <div id="downloadSelection">
                                        <div style="border: 2px solid black; height: 5%;">
                                            <div id="patName" class="heading1" style="color: #FF8A00; padding: 2px 0 3px 3px; cursor: default;" title="<%=patName%>">
                                                <%=patName%>
                                            </div>
                                        </div>

                                        <div id="previews" class="preview">
                                            <jsp:include page="downloadStudy.jsp?patient=${param.patient}&study=${param.study}&dcmURL=${param.dcmURL}&wadoUrl=${param.wadoUrl}&descDisplay=true&studyDesc=${param.studyDesc}&studyDate=${param.studyDate}&contentType=${param.contentType}"></jsp:include>

                                        </div>

                                        
                                            <div style="padding: 5px; float: left; padding-left: 3%; color: #FF8A00; font-size: 14px">
                                                <form>
                                                    <input type="radio" name="format" checked="checked" id="jpeg"> JPEG
                                                    <input type="radio" name="format" id="dicom"> DICOM
                                                    <input type="radio" name="format" id="both"> BOTH
                                                </form>
                                            </div>
                                            <div style="padding: 8px; float: left; padding-left: 3%; text-align: center; color: #FF8A00; font-size: 14px;white-space: nowrap;">
                                                <b id="No_Images" style=" font-size: 16px; padding-right: 2%"></b>
                                            </div>
                                            <div style="padding: 5px; float: right; padding-right: 2%">
                                                <button class="button" onclick="downloadBtnClick()" id="downloadBtn">Download</button>
                                                <button class="button" data-dismiss="modal" onclick="closeDownloadModal()">Cancel</button>
                                            </div>
                                        
                                    </div>
                                    <div id="modalConfirm" class="modal-confirm">

                                        <!-- Modal content -->
                                        <div class="modal-content-confirm">
                                            <div class="title">
                                                <div class="title" id="confirm-Download">CONFIRM DOWNLOAD</div>
                                            </div>
                                            <div id="confirm" style="height: 90%; width: 100%;">
                                                <div style="margin: 6px; padding-left: 3%; color: #FF8A00; font-size: 14px">
                                                    <b id="fileInfo" style=" font-size: 16px; padding-right: 2%"></b>
                                                </div>
                                                <div style="margin: 0px; float: left; padding-left: 4%; color: #FF8A00; font-size: 14px; width: 75%">
                                                    <b id="fileSize" style=" font-size: 16px; padding-right: 2%"></b>
                                                </div>
                                                <div style="margin: 0px; float: right; padding-right: 2%">
                                                    <button id="confirmDownloadBtn" class="button1" onclick="closeDownloadModal()">Download</button>
                                                    <button id="cancelBtn" class="button1" onclick="cancelDownload()">Cancel</button>
                                                </div>
                                            </div>
                                        </div>
                                        <!-- END OF modal-content-confirm -->
                                    </div>
                                </body>

                                </html>