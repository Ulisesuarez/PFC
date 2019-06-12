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
* Guruprasath R
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
                    <%@ taglib prefix="fmt" uri="http://java.sun.com/jsp/jstl/fmt"%>
                        <%@ taglib uri="http://java.sun.com/jsp/jstl/functions" prefix="fn"%>


                            <html>

                            <head>
                                <!-- <script type="text/javascript" src="js/lib/jquery.ambiance.js"></script> -->
                                <script type="text/javascript">
                                    function selectStudy(studyCheck) {
                                        var multiCheck = document.getElementsByName("multiCheck");
                                        var seriesCheck = document.getElementsByName("seriesCheck");
                                        var totalSeries = seriesCheck.length;
                                        var multiInstance = multiCheck.length;
                                        if (studyCheck.checked) {
                                            resetSelection();
                                            for (i = 0; i < multiInstance; i++) {
                                                var multiID = multiCheck[i].id;
                                                var imageInfo = multiID.split('|&');
                                                multiCheck[i].checked = true;
                                                addTotalImage(1);
                                                addSeries(imageInfo);
                                            }
                                            for (i = 0; i < totalSeries; i++) {
                                                var seriesID = seriesCheck[i].id;
                                                var imageInfo = seriesID.split('|&');
                                                seriesCheck[i].checked = true;
                                                addTotalImage(imageInfo[8]);
                                                addSeries(imageInfo);
                                            }
                                        } else {
                                            resetSelection();
                                            for (i = 0; i < multiInstance; i++) {
                                                multiCheck[i].checked = false;
                                            }
                                            for (i = 0; i < totalSeries; i++) {
                                                seriesCheck[i].checked = false;
                                            }
                                        }

                                    }

                                    function selectSeries() {
                                        var study = document.getElementById("studyCheck");
                                        var multiCheck = document.getElementsByName("multiCheck");
                                        var seriesCheck = document.getElementsByName("seriesCheck");
                                        var totalSeries = seriesCheck.length;
                                        var multiInstance = multiCheck.length;
                                        var totalChecked = 0;
                                        resetSelection();
                                        for (i = 0; i < multiInstance; i++) {
                                            if (multiCheck[i].checked == true) {
                                                totalChecked++;
                                                var multiID = multiCheck[i].id;
                                                var imageInfo = multiID.split('|&');
                                                addTotalImage(1);
                                                addSeries(imageInfo);
                                            }
                                        }
                                        for (i = 0; i < totalSeries; i++) {
                                            if (seriesCheck[i].checked == true) {
                                                totalChecked++;
                                                var seriesID = seriesCheck[i].id;
                                                var imageInfo = seriesID.split('|&');
                                                addTotalImage(imageInfo[8]);
                                                addSeries(imageInfo);
                                            }
                                        }

                                        if ((multiInstance + totalSeries) == totalChecked) {
                                            study.checked = true;
                                        } else {
                                            study.checked = false;
                                        }
                                    }
                                </script>
                            </head>

                            <body>
                                <c:if test="${param.descDisplay=='true'}">
                                    <div id='${param.study}' class="accordion1" title="${tooltip}">
                                        <input type="checkbox" name="${param.study}" id="studyCheck" onchange="selectStudy(this)"> ${param.studyDesc}
                                        <span style="float: right; padding-right: 3%">${param.studyDate}</span>
                                    </div>
                                </c:if>
                                <div>
                                    <ser:Series patientId="${param.patient}" study="${param.study}" dcmURL="${param.dcmURL}">
                                        <c:set var="middle" value="${(numberOfImages+0.5)/2}" />
                                        <fmt:formatNumber var="middle" maxFractionDigits="0" value="${middle}" />
                                        <fmt:parseNumber var="total" type="number" value="${numberOfImages}" />

                                        <script type="text/javascript">
                                            series.push({
                                                "seriesUID": '${seriesId}',
                                                "totalInstances": '${numberOfImages}',
                                                "seriesDesc": '${seriesDesc}',
                                                "modality": '${modality}',
                                                "seriesNumber": '${seriesNumber}',
                                                "seriesDate": '${seriesDate}',
                                                "bodyPart": '${bodyPart}',
                                                "studyDesc": '${param.studyDesc}',
                                                "studyDate": '${param.studyDate}'
                                            });
                                        </script>

                                        <table class="seriesTable" id="${fn:replace(seriesId, '.','_')}_table">
                                            <tbody>
                                                <img:Image patientId="${param.patient}" study="${studyId}" series="${seriesId}" dcmURL="${param.dcmURL}">
                                                    <c:choose>
                                                        <c:when test="${multiframe=='yes'}">
                                                            <tr style="cursor: default; color: #FF8A00; font-size: 12px;">
                                                                <td style="width: 5%">
                                                                    <input type="checkbox" name="multiCheck" id="${param.patient}|& ${param.studyDesc} |&${param.studyDate}|& ${param.study}|& ${seriesDesc} |&${seriesDate} |& ${seriesId} |& ${instanceNumber} |& ${numberOfFrames}" onchange="selectSeries()">
                                                                </td>
                                                                <td style="width: 16%">

                                                                    <c:choose>
                                                                        <c:when test="${param.wadoUrl == 'C-GET' || param.wadoUrl == 'C-MOVE'}">
                                                                            <img name="${instanceNumber}" id="${fn:replace(seriesId, '.','_')}_${instanceNumber}" style="${thumbSize}" src="Wado.do?dicomURL=${param.dcmURL}&study=${param.study}&series=${seriesId}&object=${imageId}&retrieveType=${param.wadoUrl}&sopClassUID=${sopClassUID}&seriesDesc=${seriesDesc}&images=${numberOfImages}&modality=${modality}"
                                                                            />
                                                                        </c:when>
                                                                        <c:otherwise>
                                                                            <img name="${instanceNumber}" id="${fn:replace(seriesId, '.','_')}_${instanceNumber}" style="${thumbSize}" src="Image.do?serverURL=${param.wadoUrl}&study=${param.study}&series=${seriesId}&object=${imageId}&sopClassUID=${sopClassUID}&seriesDesc=${seriesDesc}&images=${numberOfImages}&modality=${modality}"
                                                                            />
                                                                        </c:otherwise>
                                                                    </c:choose>
                                                                </td>
                                                                <td style="color:#FF8A00">
                                                                    <p name="${instanceNumber}"> Multiframe-${instanceNumber}</p>
                                                                </td>
                                                                <td align="right">
                                                                    <div id="totalImgs">${numberOfFrames} Frames </div>
                                                                </td>
                                                            </tr>

                                                        </c:when>

                                                        <c:otherwise>
                                                            <c:if test="${instanceNumber==1}">
                                                                <tr style="cursor: default; color: #FF8A00">
                                                                    <td style="width: 5%">
                                                                        <input type="checkbox" name="seriesCheck" id="${param.patient}|& ${param.studyDesc} |&${param.studyDate}|& ${param.study}|& ${seriesDesc} |&${seriesDate} |& ${seriesId} |& |& ${numberOfImages}" onchange="selectSeries()">
                                                                    </td>
                                                                    <td style="width: 16%">
                                                            </c:if>

                                                            <c:choose>
                                                                <c:when test="${modality == 'SR'}">
                                                                    <!-- Structured Report -->
                                                                    <img name="${instanceNumber}" id="${fn:replace(seriesId, '.','_')}_${instanceNumber}" style="height: 75px;" src="images/SR_Latest.png" imgSrc="Image.do?serverURL=${param.wadoUrl}&study=${studyId}&series=${seriesId}&object=${imageId}&sopClassUID=${sopClassUID}"
                                                                    />
                                                                </c:when>

                                                                <c:when test="${sopClassUID == '1.2.840.10008.5.1.4.1.1.104.1'}">
                                                                    <!-- Enacpsulated PDF -->
                                                                    <img name="${instanceNumber}" id="${fn:replace(seriesId, '.','_')}_${instanceNumber}" style="${thumbSize}" src="images/pdf.png" imgSrc="Image.do?serverURL=${param.wadoUrl}&study=${studyId}&series=${seriesId}&object=${imageId}&sopClassUID=${sopClassUID}"
                                                                    />
                                                                </c:when>

                                                                <c:when test="${fn:contains(sopClassUID,'1.2.840.10008.5.1.4.1.1.9')}">
                                                                    <!-- Wave Forms -->
                                                                    <img name="${instanceNumber}" id="${fn:replace(seriesId, '.','_')}_${instanceNumber}" style="${thumbSize}" src="images/pdf.png" imgSrc="Image.do?serverURL=${param.wadoUrl}&study=${studyId}&series=${seriesId}&object=${imageId}&sopClassUID=${sopClassUID}&rid=true"
                                                                    />
                                                                </c:when>

                                                                <c:when test="${sopClassUID == '1.2.840.10008.5.1.4.1.1.66'}">
                                                                    <!-- Raw Data Storage -->
                                                                    <img name="${instanceNumber}" id="${fn:replace(seriesId, '.','_')}_${instanceNumber}" style="${thumbSize}" src="images/rawdata.png" />
                                                                </c:when>

                                                                <c:otherwise>
                                                                    <c:choose>
                                                                        <c:when test="${param.wadoUrl == 'C-GET' || param.wadoUrl == 'C-MOVE'}">
                                                                            <c:choose>
                                                                                <c:when test="${(instanceNumber==1)}">
                                                                                    <img name="${instanceNumber}" id="${fn:replace(seriesId, '.','_')}_${instanceNumber}" style="${thumbSize};" src="Wado.do?dicomURL=${param.dcmURL}&study=${param.study}&series=${seriesId}&object=${imageId}&retrieveType=${param.wadoUrl}&sopClassUID=${sopClassUID}&seriesDesc=${seriesDesc}&images=${numberOfImages}&modality=${modality}"
                                                                                    />
                                                                                </c:when>

                                                                            </c:choose>
                                                                        </c:when>
                                                                        <c:otherwise>
                                                                            <c:choose>
                                                                                <c:when test="${(instanceNumber==1)}">
                                                                                    <img name="${instanceNumber}" id="${fn:replace(seriesId, '.','_')}_${instanceNumber}" style="${thumbSize}" src="Image.do?serverURL=${param.wadoUrl}&study=${param.study}&series=${seriesId}&object=${imageId}&sopClassUID=${sopClassUID}&seriesDesc=${seriesDesc}&images=${numberOfImages}&modality=${modality}&contentType=${param.contentType}"
                                                                                    />
                                                                                </c:when>
                                                                            </c:choose>
                                                                        </c:otherwise>
                                                                    </c:choose>
                                                                </c:otherwise>
                                                            </c:choose>

                                                            <c:if test="${instanceNumber==numberOfImages}">
                                                                </td>
                                                                <td style="color:#FF8A00">
                                                                    <p name="${seriesId} | ${numberOfImages}"> ${seriesDesc}</p>
                                                                </td>
                                                                <td align="right">
                                                                    <div id="totalImgs">${numberOfImages} Imgs </div>
                                                                </td>
                                                                </tr>
                                                            </c:if>
                                                        </c:otherwise>
                                                    </c:choose>
                                                </img:Image>
                                            </tbody>
                                        </table>
                                    </ser:Series>
                                </div>
                                <script type="text/javascript">
                                    var studyId = '${param.study}';
                                    if (sessionStorage[studyId] == undefined) {
                                        var study = {
                                            "studyUID": studyId,
                                            "studyDesc": "${param.studyDesc}",
                                            "studyDate": "${param.studyDate}"
                                        };
                                        var studyData = sessionStorage["${param.patient}"];

                                        if (studyData == undefined) {
                                            studyData = new Array();
                                        } else {
                                            studyData = JSON.parse(studyData);
                                        }

                                        studyData.push(study);
                                        sessionStorage["${param.patient}"] = JSON.stringify(studyData);
                                    }
                                    sessionStorage[studyId] = JSON.stringify(series);
                                </script>
                            </body>

                            </html>