
function fillAddrInput(evt) {
  sendToInput.value = evt;
  
}

function fillAccountAddrInput(evt) {
  currentAccountUI.value = evt; 
}


function jinechange(params) {
  jineJiaoyifei.value = params.value * 0.0000002;
}
//-----------

$(function () {
  //
  function showLayer(id) {
    var layer = $('#' + id),
      layerwrap = layer.find('.hw-layer-wrap');
    layer.fadeIn();
    //
    layerwrap.css({
      'margin-top': -layerwrap.outerHeight() / 2
    }); 
  }

  //
  function hideLayer() {
    $('.hw-overlay').fadeOut();
    
  }

  $('.hwLayer-ok,.hwLayer-cancel,.hwLayer-close').on('click', function () {
    hideLayer();
  });

  //
  $('.show-layer').on('click', function () {
    var layerid = $(this).data('show-layer'); 
    showLayer(layerid);
  });

  //，
  $('.hw-overlay').on('click', function (event) {
    if (event.target == this) {
      hideLayer();
    }
  });

  //ESC
  $(document).keyup(function (event) {
    if (event.keyCode == 27) {
      hideLayer();
    }
  }); 

  }); 