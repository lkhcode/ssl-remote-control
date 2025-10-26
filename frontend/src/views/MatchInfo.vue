<template>
  <div class="match-info-container">
    <!-- 左侧信息区 -->
    <div class="info-section">
      <div class="team-info">
        <div>
          <div class="team-name">{{ displayTeamName }}</div>
          <div class="badge" :class="displayTeamColor">{{ displayTeamColor.charAt(0).toUpperCase() + displayTeamColor.slice(1) }}</div>
        </div>
      </div>

      <div class="match-status">
        <div class="status-item">
          <span class="label">{{ i18n.matchPhase }}:</span>
          <span class="value">{{ matchPhase }}</span>
        </div>
        <div class="status-item">
          <span class="label">{{ i18n.currentCommand }}:</span>
          <span class="value">{{ currentCommand }}</span>
        </div>
        <div class="status-item">
          <span class="label">{{ i18n.remainingTime }}:</span>
          <span class="value">{{ remainingTime }}</span>
        </div>
      </div>
    </div>

    <!-- 右侧按钮区 -->
    <div class="buttons-section">
      <!-- Raise Challenge Flag -->
      <button 
        class="control-btn primary"
        :disabled="!canRequestChallengeFlag" 
        @click="requestChallengeFlag"
      >
        <div class="btn-title">{{ i18n.raiseChallengeFlag }}</div>
        <div class="btn-subtitle">{{ state.challengeFlagsLeft }} {{ i18n.left }}</div>
      </button>

      <!-- Fail Ballplacement -->
      <button 
        class="control-btn primary"
        :disabled="!canFailBallplacement" 
        @click="failBallplacement"
      >
        <div class="btn-title">{{ i18n.failBallplacement }}</div>
      </button>

      <!-- Emergency Stop -->
      <button 
        class="control-btn orange"
        :class="{ 'cancel': emergencyStopRequested }"
        :disabled="!canRequestEmergencyStop" 
        @click="toggleEmergencyStop"
      >
        <div class="btn-title">{{ emergencyStopRequested ? '✓ ' + i18n.cancelEmergencyStop : '! ' + i18n.emergencyStop }}</div>
        <div v-if="emergencyStopRequested" class="btn-subtitle">{{ Math.round(state.emergencyStopIn) }}s {{ i18n.left }}</div>
      </button>

      <!-- Stop Timeout / Request Timeout -->
      <button 
        v-if="canStopTimeout"
        class="control-btn danger pulse"
        :disabled="!canStopTimeout" 
        @click="stopTimeout"
      >
        <div class="btn-title">{{ i18n.stopTimeout }}</div>
        <div class="btn-subtitle">{{ Math.round(state.timeoutTimeLeft) }}s {{ i18n.left }}</div>
      </button>
      <button 
        v-else-if="timeoutRequested"
        class="control-btn warning pulse"
        :disabled="!canRequestTimeout" 
        @click="toggleTimeout"
      >
        <div class="btn-title">{{ i18n.cancelTimeout }}</div>
        <div class="btn-subtitle">{{ state.timeoutsLeft }} {{ i18n.left }} • {{ Math.round(state.timeoutTimeLeft) }}s</div>
      </button>
      <button 
        v-else
        class="control-btn warning"
        :disabled="!canRequestTimeout" 
        @click="toggleTimeout"
      >
        <div class="btn-title">{{ i18n.requestTimeout }}</div>
        <div class="btn-subtitle">{{ state.timeoutsLeft }} {{ i18n.left }}</div>
      </button>

      <!-- Change Keeper Id -->
      <button 
        class="control-btn primary"
        :disabled="!canChangeKeeperId" 
        @click="openChangeKeeperModal"
      >
        <div class="btn-title">{{ i18n.changeKeeperId }}</div>
        <div class="btn-subtitle">{{ i18n.current }}: {{ state.keeperId }}</div>
      </button>

      <!-- Robot Substitution -->
      <button 
        v-if="robotSubstitutionRequested && state.canSubstituteRobot"
        class="control-btn success pulse"
        :disabled="!canRequestRobotSubstitution" 
        @click="toggleRobotSubstitution"
      >
        <div class="btn-title">{{ i18n.finishSubstitution }}</div>
        <div v-if="state.botSubstitutionTimeLeft > 0" class="btn-subtitle">{{ Math.round(state.botSubstitutionTimeLeft) }}s {{ i18n.left }}</div>
      </button>
      <button 
        v-else-if="robotSubstitutionRequested"
        class="control-btn warning pulse"
        :disabled="!canRequestRobotSubstitution" 
        @click="toggleRobotSubstitution"
      >
        <div class="btn-title">{{ i18n.cancelSubstitution }}</div>
        <div v-if="state.botSubstitutionTimeLeft > 0" class="btn-subtitle">{{ state.botSubstitutionsLeft }} {{ i18n.left }} • {{ Math.round(state.botSubstitutionTimeLeft) }}s</div>
        <div v-else class="btn-subtitle">{{ state.botSubstitutionsLeft }} {{ i18n.left }}</div>
      </button>
      <button 
        v-else
        class="control-btn primary"
        :disabled="!canRequestRobotSubstitution" 
        @click="toggleRobotSubstitution"
      >
        <div class="btn-title">{{ i18n.requestSubstitution }}</div>
        <div class="btn-subtitle">{{ state.botSubstitutionsLeft }} {{ i18n.left }}</div>
      </button>
    </div>

    <!-- 修改守门员弹窗 -->
    <Transition name="modal-fade">
      <div v-if="showChangeKeeperModal" class="modal-overlay" @click="closeChangeKeeperModal">
        <Transition name="modal-slide">
          <div class="modal" @click.stop>
            <h3>{{ i18n.changeKeeperId }}</h3>
            <p>{{ i18n.selectKeeperNumber }} ({{ i18n.selectKeeperCurrent }}: {{ state.keeperId }})</p>
            <div class="keeper-numbers">
              <button 
                v-for="number in 16"
                :key="number - 1"
                class="number-btn"
                :disabled="number - 1 === state.keeperId"
                @click="confirmChangeKeeper(number - 1)"
              >
                {{ number - 1 }}
              </button>
            </div>
            <div class="modal-buttons">
              <button class="btn-cancel" @click="closeChangeKeeperModal">{{ i18n.cancel }}</button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>

    <!-- 确认挑战标志弹窗 -->
    <Transition name="modal-fade">
      <div v-if="showChallengeFlagModal" class="modal-overlay" @click="closeChallengeFlagModal">
        <Transition name="modal-slide">
          <div class="modal" @click.stop>
            <h3>{{ i18n.confirmChallengeFlag }}</h3>
            <p>{{ i18n.areYouSure }}</p>
            <div class="modal-buttons">
              <button class="btn-confirm" @click="confirmChallengeFlag">{{ i18n.confirm }}</button>
              <button class="btn-cancel" @click="closeChallengeFlagModal">{{ i18n.cancel }}</button>
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, onUnmounted } from 'vue';
import {
  RemoteControlRequestType,
  RemoteControlTeamStateSchema,
  RemoteControlToController_Request,
  RemoteControlToControllerSchema,
} from '../proto/ssl_gc_rcon_remotecontrol_pb';
import { Team } from '../proto/ssl_gc_common_pb';
import { ApiController } from '../services/ApiController';
import { create } from '@bufbuild/protobuf';

const api = inject<ApiController>('api');
const state = ref(create(RemoteControlTeamStateSchema, {}));

// 语言状态
const lang = ref<'zh' | 'en'>('en');

// 国际化文本
const i18n = computed(() => {
  if (lang.value === 'zh') {
    return {
      matchPhase: '比赛阶段',
      currentCommand: '当前指令',
      remainingTime: '剩余时间',
      raiseChallengeFlag: '提出异议',
      failBallplacement: '放球失败',
      emergencyStop: '紧急停止',
      cancelEmergencyStop: '取消紧急停止',
      stopTimeout: '结束Timeout阶段',
      requestTimeout: '发送Timeout请求',
      cancelTimeout: '撤销Timeout请求',
      changeKeeperId: '更改守门员车号',
      left: '剩余',
      current: '当前',
      confirmChallengeFlag: '确认异议',
      areYouSure: '你确定要提出异议吗？此操作不可撤销',
      confirm: '确认',
      cancel: '取消',
      selectKeeperNumber: '选择守门员车号',
      selectKeeperCurrent: '当前守门员车号',
      requestSubstitution: '请求机器人更换',
      finishSubstitution: '完成更换',
      cancelSubstitution: '撤销机器人更换请求',
    };
  }
  return {
    matchPhase: 'Match Phase',
    currentCommand: 'Current Command',
    remainingTime: 'Remaining Time',
    raiseChallengeFlag: 'Raise Challenge Flag',
    failBallplacement: 'Fail Ballplacement',
    emergencyStop: 'Emergency Stop',
    cancelEmergencyStop: 'Cancel Emergency Stop',
    stopTimeout: 'Stop Timeout',
    requestTimeout: 'Request Timeout',
    cancelTimeout: 'Cancel Timeout Request',
    changeKeeperId: 'Change Keeper ID',
    left: 'left',
    current: 'Current',
    confirmChallengeFlag: 'Confirm Challenge Flag',
    areYouSure: 'Are you sure you want to raise a challenge flag? This request is irrevocable.',
    confirm: 'Confirm',
    cancel: 'Cancel',
    selectKeeperNumber: 'Select keeper number',
    selectKeeperCurrent: 'Current keeper number',
    requestSubstitution: 'Request Robot Substitution',
    finishSubstitution: 'Finish Robot Substitution',
    cancelSubstitution: 'Cancel Robot Substitution Request',
  };
});

// 比赛信息
const teamName = ref('请等待比赛开始');
const teamColor = ref('blue');
const matchPhase = ref('未开始');
const currentCommand = ref('无');
const remainingTime = ref('00:00');

// 根据 state.team 判断队伍颜色和名称
const displayTeamColor = computed(() => {
  if (state.value.team === Team.YELLOW) {
    return 'yellow';
  }
  return 'blue';
});

const displayTeamName = computed(() => {
  const colorMap: Record<number, string> = {
    [Team.YELLOW]: 'Yellow Team',
    [Team.BLUE]: 'Blue Team',
    [Team.UNKNOWN]: '请等待比赛开始',
  };
  return colorMap[state.value.team] || '请等待比赛开始';
});

// 按钮可用状态
const canRequestChallengeFlag = computed(() => 
  state.value.availableRequests.includes(RemoteControlRequestType.CHALLENGE_FLAG)
);
const canFailBallplacement = computed(() => 
  state.value.availableRequests.includes(RemoteControlRequestType.FAIL_BALLPLACEMENT)
);
const canRequestEmergencyStop = computed(() => 
  state.value.availableRequests.includes(RemoteControlRequestType.EMERGENCY_STOP)
);
const canStopTimeout = computed(() => 
  state.value.availableRequests.includes(RemoteControlRequestType.STOP_TIMEOUT)
);
const canRequestTimeout = computed(() => 
  state.value.availableRequests.includes(RemoteControlRequestType.TIMEOUT)
);
const canChangeKeeperId = computed(() => 
  state.value.availableRequests.includes(RemoteControlRequestType.CHANGE_KEEPER_ID)
);
const canRequestRobotSubstitution = computed(() => 
  state.value.availableRequests.includes(RemoteControlRequestType.ROBOT_SUBSTITUTION)
);

// 请求状态
const emergencyStopRequested = computed(() => 
  state.value.activeRequests.includes(RemoteControlRequestType.EMERGENCY_STOP)
);
const timeoutRequested = computed(() => 
  state.value.activeRequests.includes(RemoteControlRequestType.TIMEOUT)
);
const robotSubstitutionRequested = computed(() => 
  state.value.activeRequests.includes(RemoteControlRequestType.ROBOT_SUBSTITUTION)
);

// 弹窗状态
const showChangeKeeperModal = ref(false);
const showChallengeFlagModal = ref(false);

// 获取比赛信息
// 不再需要外部 API 调用，所有数据从 state.team 获取

// 按钮功能
const requestChallengeFlag = () => {
  showChallengeFlagModal.value = true;
};

const confirmChallengeFlag = () => {
  api?.Send(create(RemoteControlToControllerSchema, {
    msg: {
      case: 'request',
      value: RemoteControlToController_Request.CHALLENGE_FLAG
    }
  }));
  showChallengeFlagModal.value = false;
};

const closeChallengeFlagModal = () => {
  showChallengeFlagModal.value = false;
};

const failBallplacement = () => {
  api?.Send(create(RemoteControlToControllerSchema, {
    msg: {
      case: 'request',
      value: RemoteControlToController_Request.FAIL_BALLPLACEMENT
    }
  }));
};

const toggleEmergencyStop = () => {
  api?.Send(create(RemoteControlToControllerSchema, {
    msg: {
      case: 'requestEmergencyStop',
      value: !emergencyStopRequested.value
    }
  }));
};

const stopTimeout = () => {
  api?.Send(create(RemoteControlToControllerSchema, {
    msg: {
      case: 'request',
      value: RemoteControlToController_Request.STOP_TIMEOUT
    }
  }));
};

const toggleTimeout = () => {
  api?.Send(create(RemoteControlToControllerSchema, {
    msg: {
      case: 'requestTimeout',
      value: !timeoutRequested.value
    }
  }));
};

const openChangeKeeperModal = () => {
  showChangeKeeperModal.value = true;
};

const confirmChangeKeeper = (keeperId: number) => {
  api?.Send(create(RemoteControlToControllerSchema, {
    msg: {
      case: 'desiredKeeper',
      value: keeperId
    }
  }));
  showChangeKeeperModal.value = false;
};

const closeChangeKeeperModal = () => {
  showChangeKeeperModal.value = false;
};

const toggleRobotSubstitution = () => {
  api?.Send(create(RemoteControlToControllerSchema, {
    msg: {
      case: 'requestRobotSubstitution',
      value: !robotSubstitutionRequested.value
    }
  }));
};

// 键盘事件处理
const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') {
    if (showChallengeFlagModal.value) {
      closeChallengeFlagModal();
    }
    if (showChangeKeeperModal.value) {
      closeChangeKeeperModal();
    }
  }
  // 按 L 键切换语言
  if ((event.ctrlKey || event.metaKey) && event.key === 'l') {
    event.preventDefault();
    lang.value = lang.value === 'zh' ? 'en' : 'zh';
  }
};

api?.RegisterStateConsumer((s) => {
  state.value = s;
});

api?.RegisterRefereeConsumer((referee) => {
  remainingTime.value = referee.stageTimeLeft > 0n
    ? `${Math.floor(Number(referee.stageTimeLeft / 60000000n))}:${Math.floor(Number((referee.stageTimeLeft % 60000000n) / 1000000n)).toString().padStart(2, '0')}`
    : '00:00';

  currentCommand.value = referee.command.toString();

  matchPhase.value = referee.stage.toString();

  teamName.value = referee.blueTeamOnPositiveHalf && referee.blue
    ? referee.blue.name
    : referee.yellow?.name || 'Unknown';
});

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);
});
</script>

<style scoped>
.match-info-container {
  display: flex;
  height: 100%;
  width: 100%;
  gap: 2rem;
  padding: 2rem;
  box-sizing: border-box;
  background-color: #333;
  color: #fff;
}

/* 左侧信息区 */
.info-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2rem;
}

.team-info {
  display: flex;
  align-items: flex-start;
  gap: 1.5rem;
}

.team-name {
  font-size: 2.5rem;
  font-weight: bold;
  color: #fff;
}

.badge {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  color: #fff;
  font-weight: bold;
  font-size: 1.2rem;
  margin-top: 0.5rem;
}

.badge.blue {
  background-color: #007bff;
}

.badge.yellow {
  background-color: #ffc107;
  color: #333;
}

.match-status {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.status-item {
  display: flex;
  gap: 1rem;
  font-size: 1.2rem;
}

.status-item .label {
  font-weight: bold;
  color: #aaa;
  min-width: 80px;
}

.status-item .value {
  color: #fff;
}

/* 右侧按钮区 */
.buttons-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1.5rem;
}

.control-btn {
  flex: 1;
  padding: 0.8rem 1.5rem;
  font-size: 1rem;
  font-weight: bold;
  border: none;
  border-radius: 8px;
  background-color: #007bff;
  color: #fff;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  gap: 0.2rem;
  text-align: left;
}

.btn-title {
  font-size: 1rem;
  font-weight: bold;
}

.btn-subtitle {
  font-size: 0.75rem;
  font-weight: normal;
  opacity: 0.9;
}

/* 按钮颜色变体 */
.control-btn.primary {
  background-color: #007bff;
}

.control-btn.primary:hover:not(:disabled) {
  background-color: #0056b3;
}

.control-btn.orange {
  background-color: #fd7e14;
}

.control-btn.orange:hover:not(:disabled) {
  background-color: #e66000;
}

.control-btn.orange.cancel {
  background-color: #dc3545;
}

.control-btn.orange.cancel:hover:not(:disabled) {
  background-color: #c82333;
}

.control-btn.danger {
  background-color: #dc3545;
}

.control-btn.danger:hover:not(:disabled) {
  background-color: #c82333;
}

.control-btn.warning {
  background-color: #ffc107;
  color: #333;
}

.control-btn.warning:hover:not(:disabled) {
  background-color: #e0a800;
}

.control-btn.warning.cancel {
  background-color: #ffc107;
}

.control-btn.warning.cancel:hover:not(:disabled) {
  background-color: #e0a800;
}

.control-btn.success {
  background-color: #28a745;
  color: #fff;
}

.control-btn.success:hover:not(:disabled) {
  background-color: #218838;
}

.control-btn.pulse {
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.6;
  }
}

.control-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.control-btn:disabled {
  background-color: #ccc;
  color: #666;
  cursor: not-allowed;
  opacity: 0.5;
}

/* 徽章样式 */
.badge-counter,
.badge-countdown,
.badge-info {
  display: inline-block;
  padding: 0.5rem 0.75rem;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: bold;
  white-space: nowrap;
}

.badge-counter {
  background-color: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.badge-countdown {
  background-color: rgba(255, 193, 7, 0.3);
  color: #ffc107;
  min-width: 50px;
  text-align: center;
}

.badge-info {
  background-color: rgba(0, 123, 255, 0.2);
  color: #90d0ff;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #fff;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
  min-width: 300px;
  color: #333;
}

.modal h3 {
  margin-top: 0;
  color: #333;
  font-size: 1.5rem;
}

.modal p {
  color: #666;
  margin: 1rem 0;
}

.modal input {
  width: 100%;
  padding: 0.75rem;
  margin: 1rem 0;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  box-sizing: border-box;
}

.keeper-numbers {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.75rem;
  margin: 1.5rem 0;
}

.number-btn {
  padding: 0.75rem;
  font-size: 1rem;
  font-weight: bold;
  border: 2px solid #007bff;
  border-radius: 4px;
  background-color: #fff;
  color: #007bff;
  cursor: pointer;
  transition: all 0.3s;
}

.number-btn:hover:not(:disabled) {
  background-color: #007bff;
  color: #fff;
}

.number-btn:disabled {
  background-color: #007bff;
  color: #fff;
  cursor: not-allowed;
  opacity: 0.8;
}

.modal-buttons {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-confirm,
.btn-cancel {
  flex: 1;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  font-weight: bold;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-confirm {
  background-color: #28a745;
  color: #fff;
}

.btn-confirm:hover {
  background-color: #218838;
}

.btn-cancel {
  background-color: #dc3545;
  color: #fff;
}

.btn-cancel:hover {
  background-color: #c82333;
}

@media (max-width: 768px) {
  .match-info-container {
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
  }

  .team-name {
    font-size: 1.8rem;
  }

  .status-item {
    font-size: 1rem;
  }
}

/* 过渡动画 */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-slide-enter-active,
.modal-slide-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-slide-enter-from,
.modal-slide-leave-to {
  transform: scale(0.9) translateY(-20px);
  opacity: 0;
}
</style>
