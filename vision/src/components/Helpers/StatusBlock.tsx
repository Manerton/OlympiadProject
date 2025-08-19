import { Button } from "react-bootstrap";
import { ApplicationStatus } from "../../dictionary/applicationStatus";
import { FaCheckCircle, FaQuestionCircle, FaTimesCircle } from "react-icons/fa";

interface StatusIconProps {
  status: number;
}

export const StatusIcon: React.FC<StatusIconProps> = ({ status }) => {
  switch (status) {
    case ApplicationStatus.Approved:
      return <FaCheckCircle size={20} className="text-success" />;

    case ApplicationStatus.Rejected:
      return <FaTimesCircle size={20} className="text-danger" />;

    case ApplicationStatus.Pending:
      return <FaQuestionCircle size={20} className="text-warning" />;

    default:
      return <Button variant="primary">Результаты</Button>;
  }
};
