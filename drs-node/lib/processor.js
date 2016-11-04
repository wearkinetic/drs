'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { step("next", value); }, function (err) { step("throw", err); }); } } return step("next"); }); }; }

function _toConsumableArray(arr) { if (Array.isArray(arr)) { for (var i = 0, arr2 = Array(arr.length); i < arr.length; i++) { arr2[i] = arr[i]; } return arr2; } else { return Array.from(arr); } }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var Processor = function () {
	function Processor() {
		_classCallCheck(this, Processor);

		this._handlers = {};
		this.redirect = undefined;
	}

	_createClass(Processor, [{
		key: 'on',
		value: function on(action) {
			for (var _len = arguments.length, handlers = Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
				handlers[_key - 1] = arguments[_key];
			}

			this._handlers[action] = handlers;
		}
	}, {
		key: 'process',
		value: function process(cmd, conn) {
			if (this.redirect) return this.redirect.process(cmd, conn);
			var handlers = this._handlers[cmd.action];
			if (!handlers) return undefined;
			return this._trigger.apply(this, [cmd, conn].concat(_toConsumableArray(handlers)));
		}
	}, {
		key: 'respond',
		value: function respond(cmd, conn, body) {
			var response = {
				key: cmd.key,
				body: body,
				action: 'drs.response'
			};
			conn.fire(response);
		}
	}, {
		key: '_trigger',
		value: function () {
			var _ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee(cmd, conn) {
				for (var _len2 = arguments.length, handlers = Array(_len2 > 2 ? _len2 - 2 : 0), _key2 = 2; _key2 < _len2; _key2++) {
					handlers[_key2 - 2] = arguments[_key2];
				}

				var ctx, result, _iteratorNormalCompletion, _didIteratorError, _iteratorError, _iterator, _step, h;

				return regeneratorRuntime.wrap(function _callee$(_context) {
					while (1) {
						switch (_context.prev = _context.next) {
							case 0:
								ctx = {};
								result = void 0;
								_iteratorNormalCompletion = true;
								_didIteratorError = false;
								_iteratorError = undefined;
								_context.prev = 5;
								_iterator = handlers[Symbol.iterator]();

							case 7:
								if (_iteratorNormalCompletion = (_step = _iterator.next()).done) {
									_context.next = 15;
									break;
								}

								h = _step.value;
								_context.next = 11;
								return h(cmd, conn, ctx);

							case 11:
								result = _context.sent;

							case 12:
								_iteratorNormalCompletion = true;
								_context.next = 7;
								break;

							case 15:
								_context.next = 21;
								break;

							case 17:
								_context.prev = 17;
								_context.t0 = _context['catch'](5);
								_didIteratorError = true;
								_iteratorError = _context.t0;

							case 21:
								_context.prev = 21;
								_context.prev = 22;

								if (!_iteratorNormalCompletion && _iterator.return) {
									_iterator.return();
								}

							case 24:
								_context.prev = 24;

								if (!_didIteratorError) {
									_context.next = 27;
									break;
								}

								throw _iteratorError;

							case 27:
								return _context.finish(24);

							case 28:
								return _context.finish(21);

							case 29:
								return _context.abrupt('return', result);

							case 30:
							case 'end':
								return _context.stop();
						}
					}
				}, _callee, this, [[5, 17, 21, 29], [22,, 24, 28]]);
			}));

			function _trigger(_x, _x2, _x3) {
				return _ref.apply(this, arguments);
			}

			return _trigger;
		}()
	}]);

	return Processor;
}();

exports.default = Processor;