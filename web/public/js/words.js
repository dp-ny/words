var words = function() {
  var b = {};
  var pausedDuration = undefined;
  var time = undefined;
  var stopped = true;
  var DURATION = moment.duration(2, 'minutes');
  var $body = $('body');
  var $newButton = $('.words-container .actions > .new');
  var $pauseButton = $('.words-container .actions > .pause');
  var $boardSize = $('.words-container .board-size');

  var $timer = $('.words-container .timer');
  var $table = $('table.words');
  var $tableTd = $('.words-container table.words td');
  var redFlash = '#FF1E1E';
  var defaultColor = $('.words-container').css('background-color');
  var flashCount = 10;

  var init = function() {
    var initialSize = $table.css('font-size').replace('px', '');
    $boardSize.val(initialSize);
    updateTimer();
    bindActions();
    resetGame();
  };

  var resetGame = function() {
    bindTimerUpdate();
    pausedDuration = undefined;
    var now = moment();
    time = now.add(DURATION);
    $timer.html(getTimer(DURATION.hours(), DURATION.minutes(), DURATION.seconds()));
  }

  var bindActions = function() {
    $newButton.on('click', function() {
      stopped = true;
      $.ajax({
        url: 'words/new',
        success: function(data) {
          if (!data.html) {
            alert('unable to fetch new board');
          } else {
            $table.html(data.html);
            resetGame();
            stopped = false;
          }
        },
        error: function() {
          alert('an error occurred');
        }
      });
    });
    $boardSize.change(function() {
      var fontSize = Number($boardSize.val());
      $table.css('font-size', fontSize + 'px');
      var size = fontSize + 20;
      $tableTd.css('width', size + 'px');
      $tableTd.css('height', size + 'px');
    });
    $pauseButton.on('click', function() {
      if (time === undefined) {
        // Do nothing if not started
        return;
      }
      stopped = !stopped;
      if (stopped) {
        // save current time plus pausedDuration
        var now = moment();
        pausedDuration = moment.duration(time.diff(now));
      } else {
        // start based on saved time diff
        time = moment().add(pausedDuration);
        pausedDuration = undefined;
      }
    });
  };

  var setDefaultTime = function() {
    $timer.html(getDefaultTime());
  };

  var padTime = function(t) {
    var s = '' + t;
    if (s.length < 2) {
      s = '0'+s;
    }
    return s
  };

  var defaultBody = function() {
    $body.css('background-color', defaultColor);
  }

  var flashGameOver = function(count) {
    return function() {
      if (count < 0 || stopped) {
        stopped = true;
        defaultBody();
        return;
      }
      setDefaultTime();
      if (count % 2) {
        $body.css('background-color', redFlash);
      } else {
        defaultBody();
      }
      setTimeout(flashGameOver(count-1), 500);
    }
  }

  var updateTimer = function() {
    if (stopped) {
      bindTimerUpdate();
      return;
    }
    if (!time) {
      setDefaultTime();
      return;
    }
    var now = moment();
    if (time.isBefore(now)) {
      flashGameOver(flashCount)();
      return;
    }
    var timeLeft = moment(time.diff(now));
    $timer.html(getTimer(timeLeft.hours(), timeLeft.minutes(), timeLeft.seconds()));
    bindTimerUpdate();
  };

  var bindTimerUpdate = function() {
    setTimeout(updateTimer, 100);
  };

  var getDefaultTime = function() {
    return getTimer(0, 0, 0);
  }

  var getTimer = function(hours, minutes, seconds) {
    var timer = "";
    var printMissing = false;
    if (DURATION.hours() > 0) {
      timer += padTime(hours) + ":";
      printMissing = true;
    }
    if (DURATION.minutes() > 0 || printMissing) {
      timer += padTime(minutes) + ":";
      printMissing = true;
    }
    if (DURATION.seconds() > 0 || printMissing) {
      timer += padTime(seconds);
    }
    return timer;
  };

  init();
  return b;
}();
