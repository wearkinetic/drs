'use strict';

Object.defineProperty(exports, "__esModule", {
	value: true
});

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

var _uuid = require('./uuid');

var _uuid2 = _interopRequireDefault(_uuid);

var _error = require('./error');

var _processor = require('./processor');

var _processor2 = _interopRequireDefault(_processor);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { return step("next", value); }, function (err) { return step("throw", err); }); } } return step("next"); }); }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function timeout(ms) {
	return new Promise(function (resolve) {
		setTimeout(resolve, ms);
	});
}

var Connection = function () {
	function Connection(protocol) {
		var _this = this;

		_classCallCheck(this, Connection);

		this.cache = {};
		this._protocol = protocol;
		this._pending = {};
		this._processor = new _processor2.default();
		this.on = this._processor.on.bind(this._processor);

		this._interval = setInterval(function () {
			return _this._ping();
		}, 1000);
		this._ping();
	}

	_createClass(Connection, [{
		key: '_ping',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee() {
				var result;
				return regeneratorRuntime.wrap(function _callee$(_context) {
					while (1) {
						switch (_context.prev = _context.next) {
							case 0:
								_context.next = 2;
								return this.send({
									action: 'drs.ping'
								});

							case 2:
								result = _context.sent;

								this._time = {
									server: result,
									local: Date.now()
								};

							case 4:
							case 'end':
								return _context.stop();
						}
					}
				}, _callee, this);
			}));

			return function _ping() {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: 'now',
		value: function now() {
			var _time = this._time;
			var server = _time.server;
			var local = _time.local;

			return server + (Date.now() - local);
		}
	}, {
		key: 'send',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee3(cmd) {
				var _this2 = this;

				var result;
				return regeneratorRuntime.wrap(function _callee3$(_context3) {
					while (1) {
						switch (_context3.prev = _context3.next) {
							case 0:
								if (!cmd.key) cmd.key = _uuid2.default.ascending();
								_context3.next = 3;
								return new Promise(function () {
									var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee2(resolve) {
										return regeneratorRuntime.wrap(function _callee2$(_context2) {
											while (1) {
												switch (_context2.prev = _context2.next) {
													case 0:
														_this2._pending[cmd.key] = resolve;
														_context2.next = 3;
														return _this2.fire(cmd);

													case 3:
													case 'end':
														return _context2.stop();
												}
											}
										}, _callee2, _this2);
									})),
									    _this = _this2;

									return function (_x2) {
										return ref.apply(_this, arguments);
									};
								}());

							case 3:
								result = _context3.sent;

								if (!(result.action === 'drs.error')) {
									_context3.next = 6;
									break;
								}

								throw new _error.Error(result.body);

							case 6:
								if (!(result.action === 'drs.exception')) {
									_context3.next = 8;
									break;
								}

								throw new _error.Exception(result.body);

							case 8:
								return _context3.abrupt('return', result.body);

							case 9:
							case 'end':
								return _context3.stop();
						}
					}
				}, _callee3, this);
			}));

			return function send(_x) {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: 'fire',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee4(cmd) {
				return regeneratorRuntime.wrap(function _callee4$(_context4) {
					while (1) {
						switch (_context4.prev = _context4.next) {
							case 0:
								if (!cmd.key) cmd.key = _uuid2.default.ascending();

							case 1:
								if (this._closed) {
									_context4.next = 13;
									break;
								}

								_context4.prev = 2;

								this._raw.send(this._protocol.encode(cmd));
								return _context4.abrupt('break', 13);

							case 7:
								_context4.prev = 7;
								_context4.t0 = _context4['catch'](2);

							case 9:
								_context4.next = 11;
								return timeout(1000);

							case 11:
								_context4.next = 1;
								break;

							case 13:
							case 'end':
								return _context4.stop();
						}
					}
				}, _callee4, this, [[2, 7]]);
			}));

			return function fire(_x3) {
				return ref.apply(this, arguments);
			};
		}()
	}, {
		key: 'dial',
		value: function dial(transport, host) {
			var _this3 = this;

			if (this._closed) return;
			return transport.dial(host).then(function () {
				var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee6(raw) {
					return regeneratorRuntime.wrap(function _callee6$(_context6) {
						while (1) {
							switch (_context6.prev = _context6.next) {
								case 0:
									_this3._raw = raw;
									_this3.read().then(_asyncToGenerator(regeneratorRuntime.mark(function _callee5() {
										return regeneratorRuntime.wrap(function _callee5$(_context5) {
											while (1) {
												switch (_context5.prev = _context5.next) {
													case 0:
														_context5.next = 2;
														return timeout(1000);

													case 2:
														_context5.next = 4;
														return _this3.dial(transport, host);

													case 4:
													case 'end':
														return _context5.stop();
												}
											}
										}, _callee5, _this3);
									})));

								case 2:
								case 'end':
									return _context6.stop();
							}
						}
					}, _callee6, _this3);
				})),
				    _this = _this3;

				return function (_x4) {
					return ref.apply(_this, arguments);
				};
			}()).catch(function () {
				var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee7(ex) {
					return regeneratorRuntime.wrap(function _callee7$(_context7) {
						while (1) {
							switch (_context7.prev = _context7.next) {
								case 0:
									console.log(ex);
									_context7.next = 3;
									return timeout(1000);

								case 3:
									_context7.next = 5;
									return _this3.dial(transport, host);

								case 5:
								case 'end':
									return _context7.stop();
							}
						}
					}, _callee7, _this3);
				})),
				    _this = _this3;

				return function (_x5) {
					return ref.apply(_this, arguments);
				};
			}());
		}
	}, {
		key: 'read',
		value: function read() {
			var _this4 = this;

			return new Promise(function (resolve) {
				_this4._raw.on('data', function () {
					var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee8(data) {
						var cmd, waiting, result;
						return regeneratorRuntime.wrap(function _callee8$(_context8) {
							while (1) {
								switch (_context8.prev = _context8.next) {
									case 0:
										cmd = _this4._protocol.decode(data);

										if (!(cmd.action === 'drs.response' || cmd.action === 'drs.error' || cmd.action === 'drs.exception')) {
											_context8.next = 8;
											break;
										}

										waiting = _this4._pending[cmd.key];

										if (!waiting) {
											_context8.next = 7;
											break;
										}

										waiting(cmd);
										delete _this4._pending[cmd.key];
										return _context8.abrupt('return');

									case 7:
										return _context8.abrupt('return');

									case 8:
										_context8.prev = 8;
										_context8.next = 11;
										return _this4._processor.process(cmd);

									case 11:
										result = _context8.sent;

										_this4._processor.respond(cmd, _this4, result);
										_context8.next = 18;
										break;

									case 15:
										_context8.prev = 15;
										_context8.t0 = _context8['catch'](8);

										_this4._processor.respond(cmd, _this4, _context8.t0);

									case 18:
									case 'end':
										return _context8.stop();
								}
							}
						}, _callee8, _this4, [[8, 15]]);
					})),
					    _this = _this4;

					return function (_x6) {
						return ref.apply(_this, arguments);
					};
				}());

				_this4._raw.on('close', function () {
					console.log('closed');
					resolve();
				});
			});
		}
	}, {
		key: 'close',
		value: function close() {
			var _this5 = this;

			this._closed = true;
			clearInterval(this._interval);
			Object.keys(this._pending).map(function (key) {
				_this5._pending[key].resolve({
					key: key,
					action: 'drs.exception',
					body: {
						message: 'Closing connection'
					}
				});
				delete _this5._pending[key];
			});
			return this._raw.close();
		}
	}], [{
		key: 'dial',
		value: function () {
			var ref = _asyncToGenerator(regeneratorRuntime.mark(function _callee9(transport, protocol, host) {
				var raw;
				return regeneratorRuntime.wrap(function _callee9$(_context9) {
					while (1) {
						switch (_context9.prev = _context9.next) {
							case 0:
								_context9.next = 2;
								return transport.dial(host);

							case 2:
								raw = _context9.sent;
								return _context9.abrupt('return', new Connection(raw, protocol));

							case 4:
							case 'end':
								return _context9.stop();
						}
					}
				}, _callee9, this);
			}));

			return function dial(_x7, _x8, _x9) {
				return ref.apply(this, arguments);
			};
		}()
	}]);

	return Connection;
}();

exports.default = Connection;