import type React from "react";
import type { ApplicationEvent } from "../../../../types/event";


type Props = {
    event: ApplicationEvent
    onBack: () => void
}

const ResultByEvent: React.FC<Props> = ({ event, onBack }) => {
    return (
        <div>
            <button className="btn btn-link mb-3" onClick={onBack}>
                Вернуться
            </button>
            <div className="accordion">
                <div className="accordion-item">
                    <h2 className="accordion-header" id="panelsStayOpen-headingOralPart">
                        <button className="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#panelsStayOpen-collapseOralPart" aria-expanded="true" aria-controls="panelsStayOpen-collapseOralPart">
                            Устная часть
                        </button>
                    </h2>
                    <div id="panelsStayOpen-collapseOralPart" className="accordion-collapse collapse show" aria-labelledby="panelsStayOpen-headingOralPart">
                        <div className="accordion-body">
                            <table className="table w-100">
                                <thead>
                                    <tr>
                                        <th scope="col" className="col-2">#</th>
                                        <th scope="col" className="col-6">Название</th>
                                        <th scope="col" className="col-2">Макс. Балл</th>
                                        <th scope="col" className="col-2">Полученный Балл</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
                <div className="accordion-item">
                    <h2 className="accordion-header" id="panelsStayOpen-headingPracticalPart">
                        <button className="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#panelsStayOpen-collapsePracticalPart" aria-expanded="true" aria-controls="panelsStayOpen-collapsePracticalPart">
                            Практическая часть
                        </button>
                    </h2>
                    <div id="panelsStayOpen-collapsePracticalPart" className="accordion-collapse collapse show" aria-labelledby="panelsStayOpen-headingPracticalPart">
                        <div className="accordion-body">
                            
                        </div>
                    </div>
                </div>
                <div className="accordion-item">
                    <h2 className="accordion-header" id="panelsStayOpen-headingTesting">
                        <button className="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#panelsStayOpen-collapseTesting" aria-expanded="true" aria-controls="panelsStayOpen-collapseTesting">
                            Тестирование
                        </button>
                    </h2>
                    <div id="panelsStayOpen-collapseTesting" className="accordion-collapse collapse show" aria-labelledby="panelsStayOpen-headingTesting">
                        <div className="accordion-body">
                            
                        </div>
                    </div>
                </div>  

            </div>
        </div>
    )
}

export default ResultByEvent