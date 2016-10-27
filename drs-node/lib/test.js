'use strict';

var start = function () {
	var _ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee() {
		var conn, result;
		return regeneratorRuntime.wrap(function _callee$(_context) {
			while (1) {
				switch (_context.prev = _context.next) {
					case 0:
						_context.prev = 0;
						conn = new _connection2.default(_json2.default);
						_context.next = 4;
						return conn.dial(transport, 'localhost:12000');

					case 4:
						conn.on('test', function () {
							console.log('ok');
						});
						conn.read();
						_context.next = 8;
						return conn.send({
							action: 'delta.mutation',
							body: {
								op: {
									'node.test': {
										$merge: {
											nice: Date.now()
										}
									}
								}
							}
						});

					case 8:
						result = _context.sent;

						console.log(result);
						_context.next = 16;
						break;

					case 12:
						_context.prev = 12;
						_context.t0 = _context['catch'](0);

						console.log(_context.t0);
						console.log(_context.t0.stack);

					case 16:
					case 'end':
						return _context.stop();
				}
			}
		}, _callee, this, [[0, 12]]);
	}));

	return function start() {
		return _ref.apply(this, arguments);
	};
}();

var _websocket = require('./transport/websocket');

var _websocket2 = _interopRequireDefault(_websocket);

var _json = require('./protocol/json');

var _json2 = _interopRequireDefault(_json);

var _connection = require('./connection');

var _connection2 = _interopRequireDefault(_connection);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { step("next", value); }, function (err) { step("throw", err); }); } } return step("next"); }); }; }

var transport = new _websocket2.default({
	token: 'ZGhD0ptkn0XTk7DbphrvwXPwuOPI2pEB5qIBWqjq',
	device: String(Math.random())
});

start();