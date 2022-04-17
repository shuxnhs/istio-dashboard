<template>
  <page-header-wrapper>
    <template v-slot:content>
      <div class="page-header-content">
        <div class="avatar">
          <a-avatar size="large" :src="currentUser.avatar" />
        </div>
        <div class="content">
          <div class="content-title">
            {{ timeFix }}，{{ user.name }}<span class="welcome-text">，{{ welcome }}</span>
          </div>
        </div>
      </div>
    </template>
    <template v-slot:extraContent>
      <div class="extra-content">
        <div class="stat-item">
          <a-statistic title="托管网格数量" :value="10" />
        </div>
      </div>
    </template>

    <div>
      <a-row :gutter="24">
        <a-col :xl="16" :lg="24" :md="24" :sm="24" :xs="24">
          <a-card
            class="project-list"
            :loading="loading"
            style="margin-bottom: 24px;"
            :bordered="false"
            title="服务网格"
            :body-style="{ padding: 0 }"
          >
            <div>
              <a-card-grid class="project-card-grid" :key="i" v-for="(item, i) in projects">
                <a-card :bordered="false" :body-style="{ padding: 2 }">
                  <a-card-meta>
                    <div slot="title" class="card-title">
                      <a-avatar size="small" :src="item.cover" />
                      <a>{{ item.title }}</a>
                    </div>
                    <div slot="description" class="card-description">
                      {{ item.description }}
                    </div>
                  </a-card-meta>
                </a-card>
              </a-card-grid>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </page-header-wrapper>
</template>

<script>

import { timeFix } from '@/utils/util'
import { mapState } from 'vuex'
import { PageHeaderWrapper } from '@ant-design-vue/pro-layout'
import { Radar } from '@/components'
import { listProjects } from '@/api/project'

import { getRoleList, getServiceList } from '@/api/manage'

export default {
  name: 'Workplace',
  components: {
    PageHeaderWrapper,
    Radar
  },
  data () {
    return {
      timeFix: timeFix(),
      avatar: '',
      user: {},

      projects: [],
      loading: true,
      radarLoading: true,
      activities: [],
      teams: [],

      // data
      axis1Opts: {
        dataKey: 'item',
        line: null,
        tickLine: null,
        grid: {
          lineStyle: {
            lineDash: null
          },
          hideFirstLine: false
        }
      },
      axis2Opts: {
        dataKey: 'score',
        line: null,
        tickLine: null,
        grid: {
          type: 'polygon',
          lineStyle: {
            lineDash: null
          }
        }
      },
      scale: [
        {
          dataKey: 'score',
          min: 0,
          max: 80
        }
      ],
      radarData: []
    }
  },
  computed: {
    ...mapState({
      nickname: state => state.user.nickname,
      welcome: state => state.user.welcome
    }),
    currentUser () {
      return {
        name: 'sHuXnHs',
        avatar: 'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png'
      }
    },
    userInfo () {
      return this.$store.getters.userInfo
    }
  },
  created () {
    this.user = this.userInfo
    this.avatar = this.userInfo.avatar

    getRoleList().then(res => {
      // console.log('workplace -> call getRoleList()', res)
    })

    getServiceList().then(res => {
      // console.log('workplace -> call getServiceList()', res)
    })
  },
  mounted () {
    this.getProjects()
  },
  methods: {
    getProjects () {
      this.$http.get('/list/search/projects').then(res => {
        this.projects = res.result && res.result.data
        this.loading = false
      })

      listProjects({}).then(res => {
        console.log(res.data)
      }).catch(err => {
        this.requestFailed(err)
      })
    }
  }
}
</script>

<style lang="less" scoped>
@import './Workplace.less';

.project-list {
  .card-title {
    font-size: 0;

    a {
      color: rgba(0, 0, 0, 0.85);
      margin-left: 12px;
      line-height: 24px;
      height: 24px;
      display: inline-block;
      vertical-align: top;
      font-size: 14px;

      &:hover {
        color: #1890ff;
      }
    }
  }

  .card-description {
    color: rgba(0, 0, 0, 0.45);
    height: 44px;
    line-height: 22px;
    overflow: hidden;
  }

  .project-item {
    display: flex;
    margin-top: 8px;
    overflow: hidden;
    font-size: 12px;
    height: 20px;
    line-height: 20px;

    a {
      color: rgba(0, 0, 0, 0.45);
      display: inline-block;
      flex: 1 1 0;

      &:hover {
        color: #1890ff;
      }
    }
  }

  .ant-card-meta-description {
    color: rgba(0, 0, 0, 0.45);
    height: 44px;
    line-height: 22px;
    overflow: hidden;
  }
}

.item-group {
  padding: 20px 0 8px 24px;
  font-size: 0;

  a {
    color: rgba(0, 0, 0, 0.65);
    display: inline-block;
    font-size: 14px;
    margin-bottom: 13px;
    width: 25%;
  }
}

.mobile {
  .project-list {
    .project-card-grid {
      width: 100%;
    }
  }

  .more-info {
    border: 0;
    padding-top: 16px;
    margin: 16px 0 16px;
  }

  .headerContent .title .welcome-text {
    display: none;
  }
}
</style>
