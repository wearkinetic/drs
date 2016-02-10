'use strict';

var start = function () {
	var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee2() {
		var _this2 = this;

		var ws;
		return regeneratorRuntime.wrap(function _callee2$(_context2) {
			while (1) {
				switch (_context2.prev = _context2.next) {
					case 0:
						ws = new _websocket2.default({
							// token: 'eW2q1S7noJzFwfLwapnQpY9bZP3B4ELGlFuwDZ3f',
							token: 'U7Vwcc2kA3XX4H9LLR5mxkhWxZY60RwKIaPFW96P',
							device: String(Math.random())
						});

						ws.on('mutation', function () {
							var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee(cmd) {
								return regeneratorRuntime.wrap(function _callee$(_context) {
									while (1) {
										switch (_context.prev = _context.next) {
											case 0:
												console.log(cmd);

											case 1:
											case 'end':
												return _context.stop();
										}
									}
								}, _callee, _this2);
							})),
							    _this = _this2;

							return function (_x) {
								return ref.apply(_this, arguments);
							};
						}());
						ws.router = function () {
							return 'drs.virginia.inboxtheapp.com';
						};
						_context2.next = 5;
						return ws.send({
							action: 'jarvis.event',
							body: {
								kind: 'convo.hello',
								context: {
									sender: 'test'
								}
							}
						});

					case 5:
					case 'end':
						return _context2.stop();
				}
			}
		}, _callee2, this);
	}));

	return function start() {
		return ref.apply(this, arguments);
	};
}();

var _websocket = require('./transport/websocket');

var _websocket2 = _interopRequireDefault(_websocket);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { return step("next", value); }, function (err) { return step("throw", err); }); } } return step("next"); }); }; }

start();