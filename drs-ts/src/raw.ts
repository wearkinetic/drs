abstract class Raw {
	public onData: (data: string) => void
	public onClose: () => void

	constructor() {
	}

	abstract close(): Promise<void>

	abstract send(data: string)
}

export default Raw
