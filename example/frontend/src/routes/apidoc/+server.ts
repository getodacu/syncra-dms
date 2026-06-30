// routes/+server.ts
import { ScalarApiReference } from '@scalar/sveltekit'
import type { RequestHandler } from './$types'
const render = ScalarApiReference({
  url: '/api/swagger/doc.json',
  agent: {
    disabled: true,
  },
  defaultOpenFirstTag: false,
  theme: 'fastify',
  hideClientButton: false,
  showSidebar: true,
  
})
export const GET: RequestHandler = () => {
  return render()
}
