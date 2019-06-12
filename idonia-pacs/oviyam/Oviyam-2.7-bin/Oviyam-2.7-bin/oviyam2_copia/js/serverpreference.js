$(document).ready(function() {
	$('button').button();
	$('#preferencetabs').tabs();
    $("#updateListener").click(function() {
        $.ajax({
            url: 'Listener.do',
            type: 'POST',
            data: {
                'aetitle': $('#listener_ae').val(),
                'port': $('#listener_port').val(),
                'action': 'Update'
            },
            dataType: 'text',
            success: function(msg1) {
                if (msg1.trim() == 'success') {
                    $.ambiance({
                        message: "Listener details updated and listener restarted!!!",
                        type: 'success'
                    });
                } else {
                    $.ambiance({
                        message: "Error while updating listener details!!!",
                        type: 'error'
                    });
                }
            }
        });
    });

    $('#updIoviyamCxt').click(function() {
        var iOvmCxt = $('#ioviyamCxt').val().trim();
        if (iOvmCxt.length == 0) {
            alert("iOviyam context should not be empty!!!");
            return;
        }
        if (iOvmCxt.indexOf("\/") != 0) {
            alert("iOviyam context must starts with /");
            return;
        }

        $.ajax({
            url: 'do/IOviyamContext',
            type: 'POST',
            data: {
                'iOviyamCxt': $('#ioviyamCxt').val(),
                'action': 'Update'
            },
            dataType: 'text',
            success: function(msg2) {
                if (msg2.trim() == 'success') {
                    $.ambiance({
                        message: "iOviyam context updated successfully!!!",
                        type: 'success'
                    });
                } else {
                    $.ambiance({
                        message: 'Error while updating iOviyam context!!!',
                        type: 'error'
                    });
                }
            }
        });
    });

    $('#updDownloadStudy').click(function() {
        var dwnStudy = $('#downloadStudy').val();
        if (dwnStudy == 'none') {
            dwnStudy = 'no';
        }
        $.ajax({
            url: 'DwnStudyConfig.do',
            type: 'POST',
            data: {
                'downloadStudy': dwnStudy,
                'action': 'update'
            },
            dataType: 'text',
            success: function(msg3) {
                if (msg3.trim() == 'success') {
                    $.ambiance({
                        message: "Download Study updated successfully!!!",
                        type: 'success'
                    });
                } else {
                    $.ambiance({
                        message: 'Error while updating Download Study!!!',
                        type: 'error'
                    });
                }
            }
        });
    });

    $.get("DwnStudyConfig.do", {
        'action': 'READ'
    }, function(data) {
        data = data.trim();
        $('#downloadStudy').val(data);
    }, 'text');

    $('#modalityList').keypress(function(e) {
        if (!((e.keyCode > 64 && e.keyCode < 91) || (e.keyCode > 96 && e.keyCode < 123) || e.keyCode == 47)) {
            e.preventDefault();
        }
    });

    $('#modalityList').keyup(function(e) {
        $(this).val($(this).val().toUpperCase());
    });

    $.get("overlayText.do", {
        'action': 'READ'
    }, function(data) {
        overlayText = data;
        if(overlayText.result =='Success'){
        	displayOverlay(overlayText);
        }else{
        	displayMsg("error", "Failed to read Text Overlay");
        }
        
    }, 'json');

    $('#resetOverlay').click(function() {
        $.get("overlayText.do", {
            'action': 'RESET'
        }, function(data) {
            overlayText = data;
            if(overlayText.result =='Success'){
            	displayOverlay(overlayText);
            }else{
            	displayMsg("error", "Failed to reset Text Overlay");
            }

        }, 'json');
    });

    $('#modalityDropDown').change(function() {
        if ($(this).val() == 'ALL') {
            $('#showImgLatrModality').hide();
        } else {
            $('#showImgLatrModality').show();
            $('#modalityList').val("");
            if ($(this).val() == 'SELECTED') {
                $('#modalityList').attr('placeholder', 'MG/DX');
            } else {
                $('#modalityList').attr('placeholder', '');
            }
        }
    });

    $('#imgLaterality').change(function() {
        if ($(this).attr('checked')) {
            $('#showImgLatr').show();
            $('#modalityDropDown').val('ALL');
            $('#showImgLatrModality').hide();
        } else {
            $('#showImgLatr').hide();
            $('#showImgLatrModality').hide();
            if ($(this).val() == 'SELECTED') {
                $('#modalityList').attr('placeholder', 'MG/DX');
            } else {
                $('#modalityList').attr('placeholder', '');
            }
        }
    });

    $('#updateOverlay').click(function() {
        var display = '';
        var allModality = '';
        var modalities = '';

        display = ($('#imgLaterality').attr('checked')) ? 'Yes' : 'No';

        if (display == 'Yes') {
            allModality = $('#modalityDropDown').val();
            if (allModality != 'ALL') {
                modalities = $('#modalityList').val().trim();
                if (modalities.length < 2) {
                    displayMsg('error', 'please enter correct modality List');
                    return;
                }
            }
        }
        
        var textOverlay = {
            'imageLaterality': {
               	'display': display,
               	'modality': allModality,
               	'modalityList': modalities
            }
        };

        $.ajax({
            url: 'overlayText.do',
            type: 'POST',
            data: {
                'action': 'UPDATE',
                'data': JSON.stringify(textOverlay)
            },
            dataType: 'json',
            success: function(data) {
                if (data.result.includes('Success')) {
                    displayMsg('success', 'text overlay updated successfully');
                } else {
                    displayMsg('error', 'failed to update text overlay');
                }
            }
        });
    });
});

function displayOverlay(overlayText) {
	var imgLaterality = overlayText.imageLaterality;
	
    if (imgLaterality.display.trim() == 'Yes') {
        $('#imgLaterality').prop('checked', true);
        $('#showImgLatr').show();
        $('#modalityDropDown').val(imgLaterality.modality.trim());
        if (imgLaterality.modality.trim() == 'ALL') {
            $('#showImgLatrModality').hide();
        } else {
            $('#showImgLatrModality').show();
            $('#modalityList').val(imgLaterality.modalityList.trim());
        }
    } else {
        $('#imgLaterality').prop('checked', false);
        $('#showImgLatr').hide();
        $('#showImgLatrModality').hide();
    }
}

function displayMsg(type, msg) {
    $.ambiance({
        message: msg,
        type: type
    });
}