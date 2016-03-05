"use strict";

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { return step("next", value); }, function (err) { return step("throw", err); }); } } return step("next"); }); }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var Connection = function () {
	function Connection(rw, protocol) {
		_classCallCheck(this, Connection);

		this.raw = rw;
		this._protocol = protocol;
		this._cache = {};
	}

	_createClass(Connection, [{
		key: "set",
		value: function set(key, value) {
			this._cache[key] = value;
		}
	}, {
		key: "get",
		value: function get(key) {
			return this._cache[key];
		}
	}, {
		key: "send",
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee(cmd) {
				var data;
				return regeneratorRuntime.wrap(function _callee$(_context) {
					while (1) {
						switch (_context.prev = _context.next) {
							case 0:
								_context.next = 2;
								return this._protocol.encode(cmd);

							case 2:
								data = _context.sent;

								this.raw.send(data);

							case 4:
							case "end":
								return _context.stop();
						}
					}
				}, _callee, this);
			}));

			return function send(_x) {
				return ref.apply(this, arguments);
			};
		}()
	}]);

	return Connection;
}();

exports.default = Connection;