
// Initialisation

const app = new Vue({
	el: '#app',
	data: {
		view: 'index',
		name: '<name>',
		uuid: '<uuid>',
		training_loss: '<training_loss>',
		epoch: '<epoch>',
		stats: {},
		args: {},
		images: {},
		index: {},
		dockernewcreating: false,
		dockernewerror: null,
		dockernewstatus: null,
	},
	methods: {
		playClick: () => { playClick() },
		pauseClick: () => { pauseClick() },
		stopClick: () => { stopClick() },
		backwardClick: () => { },
		forwardClick: () => { },

		aiDelete: (uuid) => {
			fetch('/api/ai/' + uuid, {
				method: 'DELETE',
			}).then(throwIfNotOk).then(o => {
			}).catch(e => {
				console.error(e)
			})
		},

		dockerNewSubmit: (e) => {
			e.preventDefault()
			app.dockernewcreating = true
			app.dockernewstatus = 'Creating container...'
			const form = new FormData(e.target)
			const data = formToJSON(form)
			fetch('/api/docker', {
				method: 'PUT',
				body: data,
			}).then(throwIfNotOk).then(o => {
				// await New event
				
				app.dockernewstatus = 'Waiting for container to start...'
			}).catch(e => {
				app.dockernewcreating = false
				app.dockernewstatus = null
				app.dockernewerror = e.toString()
			})
		},

		arrayWithAddWithSplit: (arr, n=3) => {
			let arr1 = []
			let arr2 = []
			for (let key in arr) {
				arr2.push(Object.assign({
					uuid: key,
				}, arr[key]))
				if (arr2.length >= n) {
					arr1.push(arr2)
					arr2 = []
				}
			}
			arr2.push('add')
			while (arr2.length < n) {
				arr2.push('invisible')
			}
			arr1.push(arr2)
			arr2 = []
			return arr1
		},

		formatEpoch: e => e.toString().padStart(9, '0').match(/.{1,3}/g).join(','),
		formatLoss: l => l.toPrecision(7),
	},
})

const formToJSON = form => {
	let result = {}
	for (let entry of form.entries()) {
		result[entry[0]] = entry[1]
	}
	return JSON.stringify(result)
}

const throwIfNotOk = e => {
	if (e.ok) {
		return e.json()
	}
	throw Error(e.status)
}

fetch('/api/ai').then(throwIfNotOk).then(data => {
	app.index = data
}).catch(e => {console.error(e)})

const showIndex = () => {
	app.view = 'index'
}

const showAi = uuid => {
	app.uuid = uuid
	fetch('/api/ai/' + app.uuid).then(throwIfNotOk).then(data => {
		app.name = data.name
		app.args = data.args
		app.training_loss = data.lastupdate.training_loss
		app.epoch = data.lastupdate.epoch
		app.stats = data.lastupdate.stats
		app.images = data.lastupdateimages

		app.view = 'ai'

		// tragic
		setTimeout(() => {
			graphInit()
		}, 20);
	}).catch(e => {
		console.error(e)
		window.location.hash = ''
	})
}

const showDockerNew = () => {
	app.view = 'dockernew'
}

const checkHash = () => {
	let match
	if ((match = window.location.hash.match('#!/ai/')) && match && match.index == 0) {
		showAi(window.location.hash.slice(6))
	} else if (window.location.hash == '#!/docker/new') {
		showDockerNew()
	} else {
		showIndex()
	}
}

checkHash()
window.addEventListener('hashchange', checkHash, false)

// Stream events

const stream = new EventSource('/api/stream')

stream.addEventListener('New', (e) => {
	const o = JSON.parse(e.data)
	const data = o.data
	const uuid = o.uuid
	Vue.set(app.index, uuid, data)
})

stream.addEventListener('Update', (e) => {
	const o = JSON.parse(e.data)
	const data = o.data
	const uuid = o.uuid
	if (app.index[uuid]) {
		app.index[uuid].lastupdate = data;
	}
	if (app.uuid != uuid) {
		return
	}
	app.training_loss = data.training_loss
	app.epoch = data.epoch
	app.stats = data.stats

	// tragic
	graphAddDatapoint(data.training_loss)
})

stream.addEventListener('UpdateImage', (e) => {
	const o = JSON.parse(e.data)
	const data = o.data
	const uuid = o.uuid
	if (app.index[uuid]) {
		app.index[uuid].lastupdateimage = data;
	}
	if (app.uuid != uuid) {
		return
	}
	app.images[data.id] = data
})

stream.addEventListener('Delete', (e) => {
	const o = JSON.parse(e.data)
	const data = o.data
	const uuid = o.uuid
	Vue.delete(app.index, uuid)
})

// Final steps

document.querySelector('#app').classList.remove('d-none')
