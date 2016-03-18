abstract class Raw {
	public onData: (data: string) => void
	public onClose: () => void

	constructor() {
	}

	abstract close(): Promise<void>
}

export default Raw
