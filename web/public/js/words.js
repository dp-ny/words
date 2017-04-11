var words = function() {
  var TimerPrinter = function(duration) {

  };
  /**
   * @param duration
   * @return {string} time remaining with minutes and seconds
   */
  TimerPrinter.getTimer = function(duration) {
    var pad = function(t, ignoreEmpty) {
      if (ignoreEmpty && t === 0) {
        return '';
      } else {
        return t < 10 ? '0' + t : t;
      }
    };
    if (duration.asMilliseconds() < 0) {
      duration = moment.duration(0);
    }
    var timer = pad(duration.seconds());
    timer = pad(duration.minutes()) + ":" + timer;
    tiemr = pad(duration.hours(), true) + ":" + timer;
    return timer;
  };

  // A words game object
  var Words = function(options) {
    var $container = options.$container || $('.game');
        $timer = options.$timer || $('.timer');
    var id = options.id || undefined,
        stopped = options.stopped || true,
        duration = options.duration || moment.duration(2, 'minutes'),
        time = options.time || moment.now(),
        html = options.html || '',
        gameOverFunc = options.gameOverFunc || function(){};

    var PublicWords = {};

    PublicWords.Id = function() {
      return id;
    };
    PublicWords.Stopped = function() {
      return stopped;
    };
    PublicWords.Duration = function() {
      return duration;
    };
    PublicWords.Time = function() {
      return time;
    };
    PublicWords.Html = function() {
      return html;
    };
    PublicWords.UpdateGame = function(data) {
      // TODO: consider making an object for all the fields... SeemsGood
      if (data.id !== undefined) {
        id = data.id;
      }
      if (data.stopped !== undefined) {
        stopped = data.stopped;
      }
      if (data.duration !== undefined) {
        if (moment.isMoment(data.duration)) {
          duration = data.duration;
        } else {
          duration = moment.duration(data.duration);
        }
      }
      if (data.time !== undefined) {
        if (!moment.isMoment(data.time)) {
          time = moment(data.time);
        } else {
          time = data.time;
        }
      }
      if (data.html !== undefined) {
        html = data.html;
      }
      setTimeLeft(moment());
      updateView();
    };
    PublicWords.UpdateStopped = function(s) {
      stopped = s;
      updatePause();
    };
    PublicWords.ToggleStopped = function() {
      PublicWords.UpdateStopped(!stopped);
    };
    PublicWords.UpdateSize = function(size) {
      if (size instanceof Number && !isNaN(size)) {
        // Do nothing either way
      }
    };

    /**
     * @param {string} id the id of the game
     */
    var updateId = function() {
      if (id === undefined) {
        return;
      }
      $idContainer.show();
      $id.html(id);
      $id.prop('href', 'words/' + id);
    };

    var updatePause = function() {
      if (stopped) {
        $pauseButton.text('Start');
      } else {
        $pauseButton.text('Pause');
      }
    };

    var updateHtml = function() {
      if (html !== undefined && html !== '') {
        $table.html(html);
      }
    };

    var updateView = function() {
      updateHtml();
      updateId();
      updatePause();
      updateTimer();
      // duration = moment.duration(duration);
    };

    var setDefaultTime = function() {
      showTimeLeft(moment.duration(0));
    };

    var showTimeLeft = function(duration) {
      $timer.html(TimerPrinter.getTimer(duration));
    };

    var setTimeLeft = function(now) {
      if (stopped) {
        showTimeLeft(duration);
        return;
      }
      var timeLeft = moment.duration(time.diff(now));
      showTimeLeft(timeLeft);
    };

    var updateTimer = function() {
      if (stopped) {
        setTimeLeft(now);
        bindTimerUpdate();
        return;
      }
      if (!time) {
        setDefaultTime();
        return;
      }
      var now = moment();
      if (time.isBefore(now)) {
        gameOverFunc();
        return;
      }
      setTimeLeft(now);
      bindTimerUpdate();
    };

    var bindTimerUpdate = function() {
      setTimeout(updateTimer, 100);
    };

    setDefaultTime();
    return PublicWords;
  };

  var $timer = $('.words-container .timer');
  var $container = $('.words-container ')

  var game = new Words({
    $timer: $timer,
    $container: $container
  });

  var $body = $('body');
  var $newButton = $('.words-container .actions > .new');
  var $pauseButton = $('.words-container .actions > .pause');
  var $boardSize = $('.words-container .board-size');
  var $idContainer = $('.words-container .id');
  var $id = $idContainer.find('.room-id');
  var $table = $('table.words');
  var $tableTd = $table.find('td');

  var redFlash = '#FF1E1E';
  var defaultColor = $('.words-container').css('background-color');
  var flashCount = 10;

  var init = function() {
    var initialSize = $table.css('font-size').replace('px', '');
    $boardSize.val(initialSize);
    bindActions();
    loadExistingGame();
  };

  var loadExistingGame = function() {
    var d = $timer.data('duration');
    if (d !== '') {
      if (!$.isNumeric(d)) {
        d = parseInt(d);
      }
      d = moment.duration(d);
    }
    var stopped = $timer.data('stopped');
    var t = $timer.data('time');
    if (!stopped && t !== '') {
      t = moment(t);
    }
    game.UpdateGame({
      duration: d,
      stopped: stopped,
      time: t
    });
  };

  /**
   * @param  {Object} data The html response to parse data from
   */
  var newGame = function(data) {
    game.UpdateGame(data);
    // $table.html(game.html);
    // updateId(game.id);
    // updateStopped(game.stopped);
    // setTime(moment(game.time));
    // game.duration = moment.duration(game.duration);
    // bindTimerUpdate();
  };


  /**
   * binds various actions onto the buttons and board size
   */
  var bindActions = function() {
    $newButton.on('click', function() {
      game.UpdateStopped(true);
      $.ajax({
        method: 'POST',
        url: 'words/new',
        success: function(data) {
          console.log('data', data);
          if (!data.html) {
            alert('unable to fetch new board');
          } else {
            newGame(data);
          }
        },
        error: function() {
          alert('an error occurred');
        }
      });
    });
    $boardSize.change(function() {
      var fontSize = Number($boardSize.val());
      game.UpdateSize(fontSize);
      // $table.css('font-size', fontSize + 'px');
      // var size = fontSize + 20;
      // $tableTd.css('width', size + 'px');
      // $tableTd.css('height', size + 'px');
    });
    $pauseButton.on('click', function() {
      // TODO: Maybe remove
      if (game.Time() === undefined) {
        // Do nothing if not started
        return;
      }

      game.ToggleStopped();
      $.ajax({
        method: 'POST',
        url: 'words/' + game.Id() + '/time',
        data: {stopped: game.Stopped()}
      }).done(function(data) {
        console.log('success', data);
        game.UpdateGame(data);
      })
      .fail(function(jqXHR) {
        console.log('fail', jqXHR);
      });
    });
  };

  var defaultBody = function() {
    $body.css('background-color', defaultColor);
  };

  var gameOverFunc = function() {
    flashGameOver(flashCount);
  }

  var flashGameOver = function(count) {
    return function() {
      if (count < 0 || stopped) {
        game.UpdateStopped(true);
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
  };

  init();
}();
