// components/RegionalStageList.tsx
import React from "react";

import RegionalStageCard from "./RegionalStageCard";

interface MyEvent {
    id?: string;
    name: string;
    start_date: string;
    end_date: string;
    previous_event_id?: string;
    event_type: string;
    subject?: number;
    class_number?: number;
    additional_info?: string;
    events?: MyEvent[];
}

interface Props {
  stages: MyEvent[];
  onStageClick?: (id: string) => void;
}

const RegionalStageList: React.FC<Props> = ({ stages, onStageClick }) => {
  return (
    <>
      {stages.map((stage) => (
        <RegionalStageCard
          key={stage.id}
          stage={stage}
          onClick={onStageClick}
        />
      ))}
    </>
  );
};

export default RegionalStageList;
