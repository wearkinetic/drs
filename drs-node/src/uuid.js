let lastPushTime = 0
const PUSH_CHARS = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz'
const PUSH_CHARS_REVERSE = PUSH_CHARS.split('').reverse().join('')
const lastRandChars = []


class UUID {

	now() {
		return new Date().getTime()
	}

	ascending() {
		return generate(this.now(), PUSH_CHARS)
	}

	descending() {
		return generate(this.now(), PUSH_CHARS_REVERSE)
	}

}

function generate(time, pool) {
	let len = pool.length
	let now = time
	let duplicateTime = (now === lastPushTime)
	lastPushTime = now

	let timeStampChars = new Array(8)
	for (let i = 7; i >= 0; i--) {
		timeStampChars[i] = pool.charAt(now % len)
		// NOTE: Can't use << here because javascript will convert to int and lose the upper bits.
		now = Math.floor(now / len)
	}
	if (now !== 0) throw new Error('We should have converted the entire timestamp.')

	let id = timeStampChars.join('')

	if (!duplicateTime) {
		for (let i = 0; i < 12; i++) {
			lastRandChars[i] = Math.floor(Math.random() * len)
		}
	} else {
		let i
		for (i = 11; i >= 0 && lastRandChars[i] === len - 1; i--) {
			lastRandChars[i] = 0
		}
		lastRandChars[i]++
	}
	for (let i = 0; i < 12; i++) {
		id += pool.charAt(lastRandChars[i])
	}
	if(id.length !== 20)
		throw new Error('Length should be 20.')
	return id
}

export default new UUID()
