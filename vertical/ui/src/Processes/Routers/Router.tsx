import { Route, Routes } from 'react-router-dom'

import BenefitCalculationView from '../Views/BenefitCalculationView'
import ProcessDetailView from '../Views/ProcessDetailView'
import ProcessListView from '../Views/ProcessListView'

export default function ProcessRouter() {
  return (
    <Routes>
      <Route path="/list" element={<ProcessListView />} />
      <Route path="/create" element={<BenefitCalculationView />} />
      <Route path="/:id" element={<ProcessDetailView />} />
    </Routes>
  )
}
