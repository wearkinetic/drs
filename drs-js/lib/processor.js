function Processor() {
	this._pending = {}
}

Processor.prototype.handle = function(cmd) {
	console.log(cmd)
}

module.exports = Processor
